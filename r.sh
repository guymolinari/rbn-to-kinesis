for i in {00..99}
do 
    export file=0000$i
    echo $file.csv
    aws s3 cp s3://guy-rbn-2/$file.csv /tmp
    ./csv2parquet --pad-lines /tmp/$file.csv /tmp/$file.parquet ./spots.def
    aws s3 cp /tmp/$file.parquet s3://guy-rbn-parquet
    rm /tmp/$file.*
done
