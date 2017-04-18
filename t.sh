# maybe more powerful
# for mac (sed for linux is different)
dir=`echo ${PWD##*/}`
grep "robot-controller" * -R | grep -v Godeps | awk -F: '{print $1}' | sort | uniq | xargs sed -i '' "s#robot-controller#$dir#g"
mv robot-controller.ini $dir.ini

