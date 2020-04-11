key1=`date +%s`

make create key=$key1 value=1
echo
make txread key=$key1
echo
make update key=$key1 value=2
echo
make txread key=$key1
echo

key2=`date +%s`

make rename key=$key1 newkey=$key2
echo
make txread key=$key2
echo
make txhas key=$key1
echo
make txhas key=$key2
echo
make txkeys
echo
make txkeyvalues
echo
make txcount
echo
make delete key=$key2
echo
make deleteall
echo
