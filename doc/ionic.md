# ionic编译与托管

## ionic编译

拉取代码并编译

    git clone https://github.com/ionic-team/ionic.git
    git pull origin core:core
    git checkout core
    npm install
    cd packages/core
    git npm install
    npm run build

## 压缩并上传

    cd dist
    gzip -r ./

在服务器上安装nodejs

    curl -sL https://deb.nodesource.com/setup_8.x | sudo -E bash -
    sudo apt-get install -y nodejs

    mongodb://oplogr3:HdsngClfjz2017@192.168.3.12:12003/oplogadmin