# Griffon

## Installation

To install `Griffon`, follow these steps:

1. Make sure you have [Go](https://golang.org/) installed on your machine.
2. Open a terminal and run the following command to install `Griffon`:

    ```shell
    go install github.com/bensooraj/griffon/cmd/griffon@v0.0.1
    ```

3. Once the installation is complete, you can verify it by running:

    ```shell
    griffon version
    ```

    This should display the version number of `Griffon`.

## Usage

To use `Griffon`, follow these steps:

1. Open a terminal and navigate to your project directory.
2. Run the following command to initialize `Griffon` in your project:

    ```shell
    griffon init
    ```

    This will create two files `griffon.hcl` and `startup_script.my_script.sh` in your project directory.

4. Create an SSH key pair at `~/.ssh/id_ed25519.pub` by running the following command,
    ```shell
    ssh-keygen -t ed25519 -C "<your-email-id>"
    ```
5. Run the following command to create the `Vultr` resources:

    ```shell
    export VULTR_API_KEY=<your vultr API key>
    griffon create -f griffon.hcl
    ```

    This will execute the `Griffon` command and perform the tasks specified in your `griffon.yaml` file.

## License

`Griffon` is licensed under the [MIT License](LICENSE)
