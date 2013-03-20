mqb: display .bindings easily
=============================

If you are:

1. Using Websphere MQ...
1. Using a `.bindings` file...
1. Frequently checking the `.bindings` file for queue/topic bindings...
1. Getting tired of the default `JMSAdmin` tool...

then try this instead. It's a tool which opens up a `.bindings` file of choice
and displays the entries in an easier manner. Example output is for instance:

    $ ./mqb.exe
    ACARS_UPLINK_TEST|EWMS.UPLINK_TO_PO01.QC|MQ|1208
    AEX.Publish_Aerocomponentrequest_to_ESB|AEX.AEROCOMPONENTREQUEST_TO_EEB.QL|JMS|1208
    AIRCRF_USG_TO_CROCOS|EEB.SUBSCRIBE_AIRCRAFTUSAGE_TO_CROCOS.01.QC|MQ|1208
    BAS.QL|BAS.QL|JMS|1208
    CHIP_AS.CHIPMsg|TEST.CHIP_AS_CHIPMSG.QL|MQ|1208

Using common GNU/Linux or Unix tools, the output to stdout can be parsed for prettier
display.

It's a read-only tool, so it doesn't (yet?) support updating the file itself.

This is merely a tool to help me out at work, so it may or may not be fit for other's
purposes.
