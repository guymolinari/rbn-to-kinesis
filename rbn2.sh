cd /tmp
#for year in 2019 2020
for year in 2018
do
    #for month in 01 02 03 04 05 06 07 08 09 10 11 12
    for month in 12
    do
        for day in 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31
        do
          wget -O $year$month$day.zip http://www.reversebeacon.net/raw_data/dl.php?f=$year$month$day
          jar xvf $year$month$day.zip
          aws s3 cp $year$month$day.csv s3://guy-rbn-2
          rm $year$month$day.csv
        done
    done
done
   
