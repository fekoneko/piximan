package help

import "fmt"

const configHelp = //
`Run without arguments to enter interactive mode.

> piximan config | --session-id | --rules    | --max-pending       | --reset-session |
>                | --password   | --language | --delay             | --reset-rules   |
>                |              |            | --pximg-max-pending | --reset-limits  |
>                |              |            | --pximg-delay       | --reset         |

                            Authorization options
                            ---------------------
--session-id The session ID to use for pixiv.net API autorization.
 -s          The authorization is only used when it's absolutely required, other
             requests will be made anonymously.
             You can get this ID from browser cookies on https://www.pixiv.net.
             Search for a cookie named 'PHPSESSID'.
             Do not paste the value directly in the command line as it could
             be logged in the terminal history (e.g. ~/.bash_history).
             Session ID will be encrypted and stored in ~/.piximan/session-id.

--password   The master password that can be set to encrypt the provided
 -P          session ID. If omited the password will be set to an empty string.
             Similarly to the session ID, avoid pasting the value directly.

                              Download options
                              ----------------
--size       Size (resolution) of downloaded images that will be used by default.
 -s          See 'piximan help download' for more info.

--language   Language that will be used to translate work titles and descriptions
 -L          by default. See 'piximan help download' for more info.

--rules      Path to YAML file with download rules that will be applied to every download.
 -r          May be provided multiple times. Run 'piximan help rules' for more info.

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
                             -------------------
--reset-session     Remove configured session ID.

--reset-defaults    Remove configured downloader defaults (--size and --language).

--reset-rules       Remove all configured download rules.

--reset-limits      Reset request delays and limits to default values.

--reset             Reset all configuration to default values.

                                  Examples
                                  --------
# Set session ID from X11 clipboard with no password
> piximan config --session-id $(xclip -o)

# Set session ID from shell environment variable with a password
> piximan config --session-id $PHPSESSID --password $PASSWORD

# Adjust request delays and limits to be more restrictive
> piximan config --max-pending 1 --delay 5 --pximg-max-pending 1 --pximg-delay 2

# Configure global download rules from a file
> piximan config --rules './rules.yaml'

# Configure default image size and language for downloaded works
> piximan config --size 2 --language en

# Reset all configuration
> piximan config --reset
`

func RunConfig() {
	fmt.Print(configHelp)
}
