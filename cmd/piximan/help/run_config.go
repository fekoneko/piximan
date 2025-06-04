package help

import "fmt"

const CONFIG_HELP = //
`Run without arguments to enter interactive mode.

> piximan config [--session-id ...] [--password ...]

                            Authorization options
                            ---------------------
--session-id The session ID to use for pixiv.net API autorization.
 -s          The authorization is used only when it's absolutely required, other
             requests will be made anonymously.
             You can get this ID from browser cookies on https://www.pixiv.net.
             Search for a cookie named 'PHPSESSID'.
             Do not paste the value directly in the command line as it could
             be logged in the terminal history (e.g. ~/.bash_history).
             Session ID will be encrypted and stored in ~/.piximan/session-id.

--password   The master password that can be set to encrypt the provided
 -P          session ID. If omited the password will be set to an empty string.
             Similarly to the session ID, avoid pasting the value directly.

                          Request delays and limits
                          --------------------------
--max-pending       Maximum number of concurrent requests to pixiv.net.
 -m                 Default value is 1.

--delay             Delay between eachnew request to pixiv.net in seconds.
 -d                 Default value is 2.

--pximg-max-pending Maximum number of concurrent requests to i.pximg.net.
 -M                 Default value is 5.

--pximg-delay       Delay between each new request to i.pximg.net in seconds.
 -D                 Default value is 1.

                             Reset configuration
                             --------------------
--no-session Remove the configured session ID.

--default    Reset all configuration except the session ID to default values.

                                  Examples
                                  --------
# Set session ID from X11 clipboard with no password
> piximan config --session-id $(xclip -o)

# Set session ID from shell environment variable with password
> piximan config --session-id $PHPSESSID --password $PASSWORD
`

func RunConfig() {
	fmt.Print(CONFIG_HELP)
}
