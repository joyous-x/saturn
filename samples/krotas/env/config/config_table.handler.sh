#!/bin/bash
set +x

# A: getopts是bash内建命令的， 而getopt是外部命令
#
#
# B: 在不使用getopt函数的情况下，要遍历所有的输入参数可以参考如下脚本：
#
# until [ $# -eq 0 ]
# do
#     case "$1" in
#         -h|--help) echo "-h or --help"; exit 0;;
#         -a|--along) echo "-a or --along";shift;echo "after shift $#";;
#         -b|--blong) echo "-b or --blong";
#             case "$2" in
#                 bval1) echo "  b value 1";shift 2;echo "after shift $#";;
#                 #不能判断-b后面是否会有别的参数，因此不能直接shift 2
#                 *) shift; echo "  unknow b value";
#                     if [ $# -eq 0 ];then
#                         exit 0;
#                     fi
#                     shift;;
#             esac;;
#         *) echo " unknow prop $1";shift;;
#     esac
# done

### common function ###############################

print_usage() {
    echo " usage:
    -h       print help informations
    -f path  the target file which need to be handled
    -p path  the private file which contains private informations 
    "
}

get_args(){
  while getopts "f:p:h" arg; do
    case $arg in
      f) file="$OPTARG"
         ;;
      p) priv="$OPTARG"
         ;;
      h) print_usage
         exit
         ;;
      ?) echo "-> invliad argument $arg:\"$OPTARG\", please use -h for more details"
         exit 1
         ;;
    esac
  done
}

### replace_priv ###############################

replace_priv() {
    target_file=$1
    privat_file=$2
    target_tmp=$target_file.tmp
    # clean temp file
    rm -f target_tmp; cp $target_file $target_tmp
    # replace
    # 使用 while read line 读取文件时, 如果最后一行之后没有换行符, 则read读取最后一行时遇到文件结束符EOF, 循环终止
    # , 虽然此时 $line 内存有最后一行, 但程序已经没有机会再处理, 因此可以通过以下代码来解决此问题: while read t || [[ -n ${t} ]]
    while read -r line || [[ -n ${line} ]]; do 
        if [[ ${line:0:1} == "#" ]]; then
            # annotation start with '#'
            continue
        fi 
        key=${line%:*};val=${line#*:};val=$(echo $val | awk '{gsub(/^\s+|\s+$/, "");print}');
        sed -i "s/{$key}/$val/g" $target_tmp
    done < $privat_file
    # check result
    old_line=$(cat $target_file | grep -v ^$ | wc -l)
    new_line=$(cat $target_tmp  | grep -v ^$ | wc -l)
    if [[ $old_line == $new_line ]]; then 
        mv -f $target_file.tmp $target_file
        echo "~~~ replace_priv success: target_file=$target_file privat_file=$privat_file"
    else
        echo "!!! replace_priv failed: target_file=$target_file privat_file=$privat_file"
    fi
}

### main ###############################

get_args $@
replace_priv $file $priv