key1=`date +%s`

make create key=$key1 value=1
echo
make read key=$key1
make txread key=$key1
echo
make update key=$key1 value=2
echo
make read key=$key1
make txread key=$key1
echo

key2=`date +%s`

make rename key=$key1 newkey=$key2
echo
make read key=$key2
make txread key=$key2
echo
make has key=$key1
make txhas key=$key1
echo
make has key=$key2
make txhas key=$key2
echo
make keys
make txkeys
echo
make keyvalues
make txkeyvalues
echo
make count
make txcount
echo
make delete key=$key2
echo
make deleteall
echo
make count
make txcount
echo
