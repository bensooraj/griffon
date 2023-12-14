# Griffon

## Installation

To install `Griffon`, follow these steps:

1. Make sure you have [Go](https://golang.org/) installed on your machine.
2. Open a terminal and run the following command to install `Griffon`:

    ```shell
    go install github.com/bensooraj/griffon
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

    This will create two files `griffon.yaml` and `startup_script.my_script.sh` in your project directory.

4. Run the following command to create the `Vultr` resources:

    ```shell
    griffon create -f griffon.yaml
    ```

    This will execute the `Griffon` command and perform the tasks specified in your `griffon.yaml` file.

## License

`Griffon` is licensed under the [MIT License](LICENSE)
