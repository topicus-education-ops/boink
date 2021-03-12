config {
	daysToKeep = 21
	cronTrigger = '@weekend'
}

node() {
	git.checkout { }

	dockerfile.validate { }

	def img = dockerfile.build {
		name = 'applicationscaler'
	}
	
	dockerfile.publish {
		image = img
		baseTag = false
		latestTag = false
		tags = [ "1.0.1" ]
    distribute = true
	}
}

