
MAKE="make MODE=release"

opp_featuretool enable AdaptiveScheduling
opp_featuretool disable CoalaConnector
opp_featuretool disable DecentralizedScheduling
opp_featuretool disable MQTTConnector
opp_featuretool defines > src/features/features.h

cd src && opp_makemake -f --deep -u Cmdenv -Xfeatures/coalaconnector -Xfeatures/decentralizedscheduling -Xfeatures/mqttconnector --make-so -o TaskletSimulator
$MAKE clean
$MAKE -j4
