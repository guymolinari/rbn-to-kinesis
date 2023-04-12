#for file in `aws s3 ls s3://guy-rbn-data | awk '{print $4}'`
for file in `aws s3 ls s3://guy-rbn-2 | awk '{print $4}'`
do
    basename="${file%.*}"
    echo $file
    #aws s3 cp s3://guy-rbn-data/$file /tmp
    aws s3 cp s3://guy-rbn-2/$file /tmp
    ./csv2parquet --pad-lines /tmp/$file /tmp/$basename.parquet ./spots.def
    aws s3 cp /tmp/$basename.parquet s3://guy-rbn-parquet
    rm /tmp/$basename.*
done
