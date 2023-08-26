#!/bin/bash

cnt=0
inputpath=""
if [ $# -eq 1 ]; then
    inputpath=$1
else
    echo "please specify input top20slowqStateFile file path to arg"
    exit 1
fi

# ユーザ・パスワードは違う可能性あり
# データベースは複数ある可能性あり
DATABASE_NAME=`sudo mysql -uisucon -pisucon  -e "SHOW DATABASES\G" |grep Database|grep -v information_schema|grep  -v performance_schema|grep -v mysql|grep -v sys |awk '{print $2}'`

while read -r line
do
    cnt=`expr $cnt + 1`
    echo " ##### $cnt / 20 ###############"
    # echo " ###########################"
    # echo " # $cnt / 20 "
    # echo " ###########################"
    echo ""
    echo "LINE $cnt : $line"
    sudo mysql -uisucon -pisucon $DATABASE_NAME --table -e "explain ${line:0:-2};"
    echo ""
done < $inputpath


