#!/usr/bin/env bash
#
# This script is part of the GoStats project.
#
# Sebastian Ruml <sebastian.ruml@gmail.com>, 2014.08.10

# Command line options
read -r -d '' usage <<-'EOF'
  This script syncs the data with Google BigQuery.
  
  Options:
   -d, --dir        Location of the data files (required)
   -h, --help       Shows this help text
EOF

function help {
    echo "$(basename $0) [OPTIONS]..." 1<&2

    if [ $# -gt 0 ]
    then
        echo "  ERROR: ${1}" 1<&2
        echo "" 1<&2
    fi

    echo "  ${usage}"
    exit 1
}

function main {
    date_string=`date +"%Y-%m-%d"`
    home_dir="/home/gostats/bigquery/"

    trending_repos="${data_dir}/github_trending_repos-${date_string}.json"
    echo "Syncing ${trending_repos}..."
    /home/gostats/google-cloud-sdk/bin/bq load --source_format=NEWLINE_DELIMITED_JSON github.trending ${trending_repos} ${home_dir}github_schema.json
    
    starred_repos="${data_dir}/github_most_starred-${date_string}.json"
    echo "Syncing ${starred_repos}..."
    /home/gostats/google-cloud-sdk/bin/bq load --source_format=NEWLINE_DELIMITED_JSON godoc.packages ${starred_repos} ${home_dir}github_schema.json

    godoc_packages="${data_dir}/godoc_packages-${date_string}.json"
    echo "Syncing ${godoc_packages}..."
    /home/gostats/google-cloud-sdk/bin/bq load --source_format=NEWLINE_DELIMITED_JSON godoc.packages ${godoc_packages} ${home_dir}godoc_schema.json

    echo "Done"
}

while [ $# -gt 0 ]; do
    case "$1" in
        --)
            # No more options left
            shift
            break
            ;;
        -d|--dir)
            data_dir="$2";
            shift
            ;;
        -h|--help)
            help
            ;;
    esac
    shift
done

# Validate command line options
[ -z "${data_dir}" ] && help "Setting the data files location with -d is required"

main
