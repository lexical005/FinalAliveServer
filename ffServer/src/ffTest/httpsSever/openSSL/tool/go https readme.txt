http://qingwang.blog.51cto.com/505009/160626
http://www.tuicool.com/articles/aymYbmM


=========================================================================
go的http.ListenAndServeTLS需要两个特别参数，一个是服务端的私钥文件路径，另外一个是服务端的数字证书文件路径。利用openssl工具，我们可以自己生成相关私钥和自签发的数字证书。


openssl genrsa -out server.key 2048
用于生成服务端私钥文件server.key，后面的参数2048单位是bit，是私钥的长度。
openssl生成的私钥中包含了公钥的信息


openssl rsa -in server.key -out server.key.public
我们可以根据私钥生成公钥


openssl req -new -x509 -key server.key -out server.crt -days 3650
我们也可以根据私钥直接生成自签发的数字证书


x509: certificate signed by unknown authority
默认也是要对服务端传过来的数字证书进行校验的，但客户端提示：这个证书是由不知名CA签发的。
tr := &http.Transport{ 
TLSClientConfig:    &tls.Config{InsecureSkipVerify: true}, 
} 
client := &http.Client{Transport: tr} 


====================================客户端认证服务端=====================================
但对于self-signed(自签发)证书来说，接收端并没有你这个self-CA的数字证书，也就是没有CA公钥，也就没有办法对数字证 书的签名进行验证。因此如果要编写一个可以对self-signed证书进行校验的接收端程序的话，首先我们要做的就是建立一个属于自己的 CA，用该CA签发我们的server端证书，并将该CA自身的数字证书随客户端一并发布。



建立我们自己的CA，需要生成一个CA私钥和一个CA的数字证书:
openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -key ca.key -subj "/CN=tonybai.com" -days 5000 -out ca.crt



生成server端的私钥，生成数字证书请求，并用我们的ca私钥签发server的数字证书
openssl genrsa -out server.key 2048
openssl req -new -key server.key -subj "/CN=localhost" -out server.csr
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 5000


CA: 
私钥文件 ca.key 
数字证书 ca.crt


Server: 
私钥文件 server.key 
数字证书 server.crt


CSR是Cerificate Signing Request的英文缩写，即证书请求文件，也就是证书申请者在申请数字证书时由CSP(加密服务提供者)在生成私钥的同时也生成证书请求文件，证书申请者只要把CSR文件提交给证书颁发机构后，证书颁发机构使用其根证书私钥签名就生成了证书公钥文件，也就是颁发给用户的证书。



client端需要验证server端的数字证书，因此client端需要预先加载ca.crt，以用于服务端数字证书的校验
pool := x509.NewCertPool() 
caCertPath := "ca.crt"

caCrt, err := ioutil.ReadFile(caCertPath) 
if err != nil { 
fmt.Println("ReadFile err:", err) 
return 
} 
pool.AppendCertsFromPEM(caCrt)

tr := &http.Transport{ 
TLSClientConfig: &tls.Config{RootCAs: pool}, 
} 
client := &http.Client{Transport: tr} 



====================================服务端认证客户端=====================================
服务端要对客户端数字证书进行校验，首先客户端需要先有自己的证书



生成客户端的私钥与证书
openssl genrsa -out client.key 2048
openssl req -new -key client.key -subj "/CN=tonybai_cn" -out client.csr
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 5000



服务端认证
pool := x509.NewCertPool() 
caCertPath := "ca.crt"

caCrt, err := ioutil.ReadFile(caCertPath) 
if err != nil { 
fmt.Println("ReadFile err:", err) 
return 
} 
pool.AppendCertsFromPEM(caCrt)

s := &http.Server{ 
Addr:    ":8081", 
Handler: &myhandler{}, 
TLSConfig: &tls.Config{ 
ClientCAs:  pool, 
ClientAuth: tls.RequireAndVerifyClientCert , 
}, 
}

err = s.ListenAndServeTLS("server.crt", "server.key") 


客户端
pool := x509.NewCertPool() 
caCertPath := "ca.crt"

caCrt, err := ioutil.ReadFile(caCertPath) 
if err != nil { 
fmt.Println("ReadFile err:", err) 
return 
} 
pool.AppendCertsFromPEM(caCrt)

cliCrt, err := tls.LoadX509KeyPair("client.crt", "client.key") 
if err != nil { 
fmt.Println("Loadx509keypair err:", err) 
return 
}

tr := &http.Transport{ 
TLSClientConfig: &tls.Config{ 
RootCAs:      pool, 
Certificates: []tls.Certificate{cliCrt}, 
}, 
} 
client := &http.Client{Transport: tr} 


出错。从server端的错误日志来看，似乎是client端的client.crt文件不满足某些条件。
$go run server.go 
2015/04/30 22:13:33 http: TLS handshake error from 127.0.0.1:53542: 
tls: client's certificate's extended key usage doesn't permit it to be 
used for client authentication

$go run client.go 
Get error: Get https://localhost:8081: remote error: handshake failure


错误出自crypto/tls/handshake_server.go。
k := false 
for _, ku := range certs[0].ExtKeyUsage { 
if ku == x509.ExtKeyUsageClientAuth { 
ok = true 
break 
} 
} 
if !ok { 
c.sendAlert(alertHandshakeFailure) 
return nil, errors.New("tls: client's certificate's extended key usage doesn't permit it to be used for client authentication") 
}

大致判断是证书中的ExtKeyUsage信息应该包含clientAuth。翻看openssl的相关资料，了解到自CA签名的数字证书中包含的都是一些basic的信息，根本没有ExtKeyUsage的信息


查看一下当前client.crt的内容
openssl x509 -text -in client.crt -noout 


golang的tls要校验ExtKeyUsage，如此我们需要重新生成client.crt，并在生成时指定extKeyUsage


创建文件client.ext，文件内容
extendedKeyUsage=clientAuth


重建client.crt
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -extfile client.ext -out client.crt -days 5000



再次测试，ok


=========================================================================
(1)openssl genrsa -out rootCA.key 2048 
(2)openssl req -x509 -new -nodes -key rootCA.key -subj "/CN=*.tunnel.tonybai.com" -days 5000 -out rootCA.pem

(3)openssl genrsa -out device.key 2048 
(4)openssl req -new -key device.key -subj "/CN=*.tunnel.tonybai.com" -out device.csr
(5)openssl x509 -req -in device.csr -CA rootCA.pem -CAkey rootCA.key -CAcreateserial -out device.crt -days 5000

(6)cp rootCA.pem assets/client/tls/ngrokroot.crt 
(7)cp device.crt assets/server/tls/snakeoil.crt 
(8)cp device.key assets/server/tls/snakeoil.key

自己搭建ngrok服务，客户端要验证服务端证书，我们需要自己做CA，因此步骤(1)和步骤(2)就是生成CA自己的相关信息。 
步骤(1) ，生成CA自己的私钥 rootCA.key 
步骤(2)，根据CA自己的私钥生成自签发的数字证书，该证书里包含CA自己的公钥。

步骤(3)~(5)是用来生成ngrok服务端的私钥和数字证书（由自CA签发）。 
步骤(3)，生成ngrok服务端私钥。 
步骤(4)，生成Certificate Sign Request，CSR，证书签名请求。 
步骤(5)，自CA用自己的CA私钥对服务端提交的csr进行签名处理，得到服务端的数字证书device.crt。

步骤(6)，将自CA的数字证书同客户端一并发布，用于客户端对服务端的数字证书进行校验。 
步骤(7)和步骤(8)，将服务端的数字证书和私钥同服务端一并发布。