domaintest-server

## 编译
    需要完整的go环境
    ```
    make 
    make docker
    ```
## 运行
    配置选择
    - config.yaml
    - 设置环境变量

### 示例1
    ```bash 
    docker run -d -v /file/path/for/you/config.yaml:/config.yaml -p 6666:6666 domaintest-server:0.0.1
    ```
    
### 示例2
	需要先下载 ipipfree.ipdb文件
	修改docker-compose.yaml 中的volumes中的ipipfree.ipdb路径ee
    ```bash 
    MYSQL_DATABASE_PASSWORD=123456 docker-compose up
    ```
### 配置config.yaml 
    参照`config.yaml.sample`文件
    需要去掉`.sample`
    示例启动
    
### 环境变量
    將config.yaml文件的项扁平化即可

#### 示例
    ````
    docker run -d -p 6666:6666 \
    -e DATABASE_USER="root" \
    -e DATABASE_PASSWD="123456" \
    -e DATABASE_NAME="gogin" \
    -e DATABASE_PORT="3306" \
    -e DATABASE_HOST="127.0.0.1" \
    -e http_fail_threshold=0.74 \
    -v /path/you/ipip_file_path:/ipipfree.ipdb \
    domaintest-server:0.0.1
    ```
