griffon {
    region = "us-east-1"
    vultr_api_key = "VULTR_API_KEY"
}

data "region" "current" {}

data "plan" "all" {
    filter {
        region = region.current.id
        vcpu_count = 1
        ram = 1024
        disk = 20
    }
}

data "os" "centos" {
    filter {
        type = "vc2"
        name   = "CentOS 7 x64"
        arch   = "x64"
        family   = "centos"
    }
}

ssh_key "my_key" {
    ssh_key = "ssh-rsa AAAAB3NzaC1yc2E"
}

startup_script "my_script" {
    script = "#!/bin/bash\necho 'hello world'"
}

instance "my_vps" {
    region = data.region.current.id
    plan = data.plan.all.id
    os_id = data.os.centos.id

    sshkey_id = ssh_key.my_key.id
    script_id = startup_script.my_script.id

    hostname = "ben-vps"

    tags = {
        name = "ben-vps"
        tier = "web"
        env = "dev"
    }
}