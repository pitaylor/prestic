# Pete's Restic

Lets you define and run restic commands from a YAML file.

Features:
* YAML allows reuse of repository definitions and command line flags with YAML anchors and aliases
* Automatically set `--parent` flag to snapshot ID of the last backup
