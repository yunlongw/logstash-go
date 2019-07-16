<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2019/1/22/022
 * Time: 13:17
 */

namespace app\common\traits;


/**
 * Trait UdpLogTrait
 * @package App\Traits
 */
trait UdpLogTrait
{
    private $fp;
    private $server;
    private $port;
    private $stack;
    protected $defaultLogData = 'GPCS';


    /**
     * @param $log
     * @param string $name
     * @param bool $showDefaultLogData
     */
    function write($log, $name = "", $showDefaultLogData = true)
    {
        if (!$showDefaultLogData){
            $this->defaultLogData = false;
        }
        $this->connect();
        if (is_resource($this->fp)) {
            $log = $this->formatRemoteLog($log, $name);
            fwrite($this->fp, $log);
            fclose($this->fp);
        }
    }

    /**
     *
     */
    public function connect()
    {

        $this->server = config('udplog.UDP_LOG_SERVER') ?? "127.0.0.1";
        $this->port = config('udplog.UDP_LOG_PORT') ?? "9091";
        $fp = fsockopen("udp://{$this->server}", $this->port, $error_no, $error_string);
        if (!$fp) {
            return;
        }

        $this->fp = $fp;
    }

    /**
     * @param $log
     * @param $name
     * @return string
     */
    public function formatRemoteLog($log, $name)
    {
        if (is_scalar($log)) {
            $log = ['string_message' => (string)$log] ;
        } elseif (is_array($log)) {
        } elseif (is_object($log)) {
            $log = var_export($log, true);
        } else {
            $log = '!!!不支持的LOG格式!!!';
        }

        if ($name){
            $log = [
                'name' => $name,
                'current_time' => date('Y-m-d H:i:s'),
                'data' => $log,
            ];
        }else{
            $log = [
                'name' => 'default_tag',
                'current_time' => date('Y-m-d H:i:s'),
                'data' => $log,
            ];
        }

        $this->addToLog('data', $log);
        $this->getLogContent();

        return json_encode($this->stack);
    }


    function addToLog($key, $content)
    {
        $this->stack[$key] = json_encode($content);
        return $this;
    }


    function getLogContent()
    {
        if ($this->defaultLogData) {
            $tokens = str_split($this->defaultLogData);
            $allowToken = array('G' => true, 'P' => true, 'C' => true, 'S' => true);
            foreach ($tokens as $t) {
                if (isset($allowToken[$t])) {
                    switch ($t) {
                        case 'G':
                            $this->addToLog('get', $_GET);
                            break;
                        case 'P':
                            $this->addToLog('post', $_POST);
                            break;
                        case 'C':
                            $this->addToLog('cookie', $_COOKIE);
                            break;
                        case 'S':
                            $session = array();
                            if (isset($_SESSION)) {
                                $session = &$_SESSION;
                            }
                            $this->addToLog('session', $session);
                            break;
                    }
                }
            }
        }
    }
}
