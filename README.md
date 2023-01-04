# DKG-Cli

This tool is a **WORK IN PROGRESS!**, do not use in production!  

DKG-cli enables conducting a DKG ceremony for [ssv.network](https://ssv.network) over any communication channel (email, telegram, web, discord, etc).

## How it works
DKG ceremonies are message driven with pre- and post-message local calculation by each node. Those messages have public and secret filed in them. Secret fields are encrypted which makes sending the messages over any channel possible.

The ceremony is async and has 3 phases:
1) A ceremony conductor choose N nodes from the SSV network to participate, contacts them and asks them to run the 'generate-node' command. Each node them publishes the node file to the others on any medium chosen. 
2) Once all nodes submitted their node files, each node independently executes the 'deal' command and publishes the resulting file on the communication medium.
3) Once all nodes submitted their deal files, each node runs the 'finalize' command and publishes the final [Output](./dkg/output.go) struct file to everyone else.

At this point any of the participants can reconstruct the validator public key and signed deposit data. Then the validator can be registered to the SSV network like any other validator. 