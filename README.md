# iptool
    COMMAND DOCUMENTATION
    
        NAME:
    
    	iptool - Opinionated tool to perform common IP queries on connected hosts
    
        USAGE:
    
    	iptool [global options] command [command options] [arguments...]
    
        VERSION:
    
    	1.0.3
    
        AUTHOR(S):
    
    	Andres Villarroel <andres.via@gmail.com>
    
        COMMANDS:
    
    	router	Do a DNS request to myip.opendns.com to get your router IP address
    	ip		Creates a simple UDP/53 connection to Google or OpenDNS and returns the source IP address
    	lan		alias of 'ip' command
    	docker	Obtain address from $DOCKER_HOST, docker.local or local.docker. If /var/run/docker.sock is a socket uses 'lan' or 127.0.0.1 if something answers on 2375 or 2376
    	version	Version number only
    	help, h	Shows a list of commands or help for one command
    
        GLOBAL OPTIONS:
    
    	--help, -h		show help
    	--version, -v	print the version
    
    
