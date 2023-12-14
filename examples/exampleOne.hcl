griffon {
    region = "ams"
    vultr_api_key = env.VULTR_API_KEY
}

data "region" "current" {}

data "plan" "vhf_32gb" {
    filter {
        type = "all"
        region = data.region.current.id
        vcpu_count = 8
        ram = 32768
        disk = 512
    }
    depends_on = [data.region.current]
}

data "os" "centos_7" {
    filter {
        type = "vc2"
        name   = "CentOS 7 x64"
        arch   = "x64"
        family   = "centos"
    }
}

ssh_key "my_key" {
    ssh_key = "ssh-rsa AAAAB3NzaC1yc2E"
    depends_on = [data.plan.vhf_32gb]
}

startup_script "my_script" {
    script = file("example/startup_script.my_script.sh")
    depends_on = [ssh_key.my_key, data.region.current]
}

instance "my_vps" {
    region = data.region.current.id
    plan = data.plan.vhf_32gb.id
    os_id = data.os.centos_7.id

    sshkey_id = ssh_key.my_key.id
    script_id = startup_script.my_script.id

    hostname = "ben-vps"

    tags = {
        name = "ben-vps"
        tier = "web"
        env = "dev"
    }

    depends_on = [data.region.current, data.plan.vhf_32gb, data.os.centos_7, ssh_key.my_key, startup_script.my_script]
}