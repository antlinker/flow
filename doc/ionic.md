# ionic编译与托管

## ionic编译

拉取代码并编译

    git clone https://github.com/ionic-team/ionic.git
    git pull origin core:core
    git checkout core
    npm install
    cd packages/core
    npm install
    npm run build

## 压缩并上传

    cd dist
    gzip -r ./
    scp -r * ubuntu@119.27.161.206:/home/ubuntu/www

## 开始服务
在服务器上安装nodejs

    curl -sL https://deb.nodesource.com/setup_8.x | sudo -E bash -
    sudo apt-get install -y nodejs


安装http-server

    npm install http-server -g

运行http-server

    cd /home/ubuntu/www
    nohup http-server -g --cors='*' >> http-server.log &

其中cors允许跨域