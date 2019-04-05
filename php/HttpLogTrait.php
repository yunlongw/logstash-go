<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2019/2/6/006
 * Time: 12:39
 */

namespace app\common\traits;


trait HttpLogTrait
{
    private $station_server;
    private $port;
    private $timeout = 1;
    private $fp;
    private $defaultLogData = 'GPCS';
    private $stack;


    /**
     *
     * HttpLog constructor.
     */
    function connect()
    {
        $this->station_server = env('HTTP_LOG_SERVER') ?? "127.0.0.1";
        $this->port = env('HTTP_LOG_PORT') ?? "9091";

        $fp = fsockopen($this->station_server, $this->port, $error_no, $error_string, $this->timeout);
        if (!$fp) {
            return;
        }
        if (!stream_set_blocking($fp, 0)) {
            return;
        }
        $this->fp = $fp;
    }

    /**
     * 写入日志
     *
     * @param string|array $log
     * @param string $name
     * @return mixed|void
     */
    function write($log, $name = '')
    {
        $this->connect();
        if (is_resource($this->fp)) {
            $log = $this->formatRemoteLog($log, $name);
            $content_length = strlen($log);
            $q = array(
                'POST /write HTTP/1.1',
                "Host: {$this->station_server}",
                "User-Agent: LogStation Client",
                "Content-Length: {$content_length}",
                "Connection: Close\r\n",
                $log
            );

            $string = implode("\r\n", $q);
            fwrite($this->fp, $string, 40960);
            fclose($this->fp);
        }
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


    /**
     * @param $key
     * @param $content
     * @return $this
     */
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
