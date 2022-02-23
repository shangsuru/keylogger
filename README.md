# Keylogger

## Server

Start the server by running `docker-compose up -d`. The server listens at port 2345 for the keystrokes sent by the kernel module.  
It also exposes a postgres instance at port 5423 and a RabbitMQ management console at port 8080.  
An API to query the keystrokes stored in the database is available at port 5000.

## Kernel Module

Follow these steps to install the kernel module onto the linux machine.
It will capture keystrokes and send them to a hard coded (!) server address at port 2345.  
For testing, you can start a TCP listener at that server with `nc -l 2345`.

1. Compile the kernel module

    ```
    make
    ```

2. Load the kernel module

    ```
    sudo insmod keylogger.ko
    sudo lsmod # show kernel modules
    ```

3. View kernel logs

    ```
    sudo dmesg
    ```

4. Unload the kernel module
    ```
    sudo rmmod keylogger
    ```
