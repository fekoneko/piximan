package help

import "fmt"

const CONFIG_HELP = //
`Run without arguments to enter interactive mode.

> piximanctl config [--sessionid ...] [--password ...]

                            Authorization options
                            ---------------------
--sessionid  The session ID to use for pixiv.net API autorization.
 -s          The authorization is used only when it's absolutely required, other
             requests will be made anonymously.
             You can get this ID from browser cookies on https://www.pixiv.net.
             Search for a cookie named 'PHPSESSID'.
             Do not paste the value directly in the command line as it could
             be logged in the terminal history (e.g. ~/.bash_history).
             Session ID will be encrypted and stored in ~/.piximan/sessionid.
             You can remove the session ID by providing an empty string.

--password   The master password that can be set to encrypt the provided
 -p          session ID. If omited the password will be set to an empty string.
             Similarly to the session ID, avoid pasting the value directly.

                                  Examples
                                  --------
# Set session ID from X11 clipboard with no password
> piximanctl config --sessionid $(xclip -o)

# Set session ID from shell environment variable with password
> piximanctl config --sessionid $PHPSESSID --password $PASSWORD
`

func RunConfig() {
	fmt.Print(CONFIG_HELP)
}
