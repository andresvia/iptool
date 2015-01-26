# iptool

[ ![Download](https://api.bintray.com/packages/andresvia/tools/iptool/images/download.svg) ](https://bintray.com/andresvia/tools/iptool/_latestVersion)

	NAME:
	   iptool - Opinionated tool to perform common queries on connected hosts
	
	USAGE:
	   iptool [global options] command [command options] [arguments...]
	
	VERSION:
	   1.0.0
	
	AUTHOR:
	  Andres Villarroel - <andres.via@gmail.com>
	
	COMMANDS:
	   router	Do a DNS request to myip.opendns.com to get your router IP address, using the same technique as in the command `dig +short myip.opendns.com @208.67.222.222` but using GO code.
	   ip		Creates a simple UDP connection to Google or OpenDNS DNS servers and returns the source IP address
	   help, h	Shows a list of commands or help for one command
	   
	GLOBAL OPTIONS:
	   --help, -h		show help
	   --version, -v	print the version
	   
