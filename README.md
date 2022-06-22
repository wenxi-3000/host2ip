# host2ip
使用Golang协程池批量将域名转为ip

# 使用方法
host2ip -f host.txt -o ips.txt -t 100 (默认是100协程，默认保存到当前目录的ips.txt文件，默认读取当前目录的host.txt)
