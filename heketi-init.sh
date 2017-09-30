#!/bin/bash
# Create the topology for heketi
# README FIRST: Used only in k8s cluster initialized phase
# Environment as below:
# * VOLUME_PATH: path to save topology and loaded endpoint file, default to /etc/heketi. (Optional)
# * ENDPOINTS: a string like '10.10.1.0,10.10.1.1,10.10.1.2;10.20.1.0,10.20.1.1,10.20.1.2;'
# * PATHS: a string like '/disk/sata00,/disk/sata01'

sleep 2
if [ ! -z $VOLUME_PATH ]; then
    volume_path="$VOLUME_PATH"
    record_path="${VOLUME_PATH}/loaded_records"
else
    volume_path="/usr/local/heketi"
    record_path="/usr/local/heketi/loaded_records"
fi
mkdir -p $(dirname $record_path)
if [ ! -f $record_path ]; then
    touch $record_path
fi

IFS=';' read -ra endpoints <<< "$ENDPOINTS"
for endpoint in "${endpoints[@]}"; do
    if grep -q "$endpoint" $record_path; then
        echo "$endpoint already loaded; skip"
        continue
    else
        echo "loading $endpoint"
    fi

    printf "{\"clusters\": [\n\t{\"nodes\": [" > $volume_path/tmp-topology.json
    IFS=',' read -ra nodes <<< "$endpoint"
    nflag=0
    for node in "${nodes[@]}"; do
        echo "adding $node"
        if [ $nflag -eq 1 ]; then
            printf ",\n" >> $volume_path/tmp-topology.json
        else
            printf "\n" >> $volume_path/tmp-topology.json
            nflag=1
        fi
        printf "\t\t{\n\t\t\t\"node\": { \"hostnames\": { " >> $volume_path/tmp-topology.json
        printf "\"manage\": [\"$node\"], " >> $volume_path/tmp-topology.json
        printf "\"storage\": [\"$node\"]}, " >> $volume_path/tmp-topology.json
        printf "\"zone\": 1}, " >> $volume_path/tmp-topology.json
        printf "\n\t\t\t\"devices\": [" >> $volume_path/tmp-topology.json

        IFS=',' read -ra paths <<< "$PATHS"
        pflag=0
        for path in "${paths[@]}"; do
            if [ $pflag -eq 1 ]; then
                printf ",\n" >> $volume_path/tmp-topology.json
            else
                printf "\n" >> $volume_path/tmp-topology.json
                pflag=1
            fi
            printf "\t\t\t\t\"$path\"" >> $volume_path/tmp-topology.json
        done
        printf "\n\t\t\t]\n\t\t}" >> $volume_path/tmp-topology.json
    done
    printf "\n\t]}\n\t]\n}\n" >> $volume_path/tmp-topology.json

    echo -e "loaded topology as below:\n\n"
    cat $volume_path/tmp-topology.json
    echo -e "\n\n"

    heketi-cli topology load --json=$volume_path/tmp-topology.json
    echo "$endpoint" >> $record_path
done

