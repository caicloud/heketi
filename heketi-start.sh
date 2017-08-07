#!/bin/sh
# Author: Delweng Zheng <delweng@gmail.com>

usage() {
    echo "Usage: $0 -c /path/to/config/file -b /path/to/backup/dir -d /path/to/heketi/database/dbfile"
    exit 1
}


TEMP=$(getopt -o c:b:d:h --long config:,backup:,datadb:,help -- "$@")
eval set -- "$TEMP"

while true
do
    case $1 in
        -c|--config)
            config=$2; shift 2;;
        -b|--backup)
            backup=$2; shift 2;;
        -d|--datadb)
            datadb=$2; shift 2;;
        --)
            shift; break;;
        -h|--help)
            usage;;
    esac
done

config=${config-/etc/heketi/heketi.json}
backup=${backup-/backupdb/}
datadb=${datadb-/var/lib/heketi/heketi.db}
datadb_dir=$(dirname $datadb)

mkdir -p $backup
if [ ! -f $datadb ]; then
    recent=$(ls $backup | sort -r | head -1)
    if [ "$recent" != "" ]; then
        tar zxf $backup/$recent -C $datadb_dir
        ret=$?
        if [ $ret -ne 0 ]; then
            echo "Unable to recovery database($ret)"
            exit $ret
        fi
        echo "Recovery database from $backup/$recent to $datadb"
    else
        echo "Heketi fresh starting!"
    fi
fi

echo "30 2 * * * root (cd $datadb_dir; export backup_file=heketi-db-\$(date +\\%F-\\%H-\\%M).tar.gz && tar zcf \$backup_file heketi.db && mv \$backup_file $backup) >/dev/null 2>&1\n#" > /etc/crontab

pidof cron
if [ $? != 0 ]; then
    cron -L15
fi

echo "Run as '/usr/bin/heketi --config=$config --backup=$backup --datadb=$datadb'"
exec "/usr/bin/heketi --config=$config"
