#  ntfy

[ntfy](https://ntfy.sh) is a simple pub/sub server for sending notifications to your devices. If you are connecting to the public instance, I suggest using a random string for the topic name as well as enabling encryption. If you are using a self-hosted instance, I suggest creating a new user solely for clipshift that is separate from any other user accounts.

## clipshift config
Use `clipshift config add-backend` to be guided through the configuration options or you can manually edit `~/.clipshift/config.yaml` and add the following:

- Using the public server

    ```yaml
    backends:
      - type: ntfy
        host: https://ntfy.sh
        topic: some-random-string-here
        action: sync
        encryptionkey: some-key-here # optional, but recommended for the public instance
    ```

- Using a [self-hosted](https://docs.ntfy.sh/install/) instance

    ```yaml
    backends:
      - type: ntfy
        host: https://ntfy.example.com
        user: some-user
        pass: some-password
        topic: clipshift
        action: sync
        encryptionkey: # optional, leave blank for no encryption
    ```


## Android

*Note: encryption is currently not supported*

Google, instead of making clipboard monitoring access a permission, has made it super locked down. So, unless Google changes course, there will probably never be an official clipshift app. But clipshift still works great with the offical ntfy app ([Play Store](https://play.google.com/store/apps/details?id=io.heckel.ntfy), [F-Droid](https://f-droid.org/en/packages/io.heckel.ntfy/)) and [Tasker](https://tasker.joaoapps.com/).

Install both apps and subscribe to the topic you configured in your clipshift client(s) in the ntfy app.

### Ntfy reactor Tasker profile

This will be triggered when a ntfy message is receive and will set the clipboard accordingly.

1. Open Tasker, if this is your first time using the app make sure you enable full Tasker and not TaskerNet (or whatever the beginner mode is called)
1. Create a profile and name it whatever you want (I named mine `ntfy reactor`)
1. For the trigger select `Event` > `System` > `Intent Received`
1. In the `Action` field put `io.heckel.ntfy.MESSAGE_RECEIVED`
1. Set both `Cat` fields to `Default`
1. Select the back arrow, and then `New Task` and you don't have to give it a name
1. In the Task Edit screen, hit the action button to add a step, select `Task` > `If`
1. For the Condition, make the left box `%topic` and the right box `clipshift` (change accordingly if you're using a different topic name)
1. Hit the `+` button to add a second condition and make the fields `%title` and `%DEVMOD` (device model), then hit the button with the `~` for that condition and change it to `Doesn't Match` which will change the button to `!~`
1. Hit the back arrow to save the step
1. Select the action button to add a step and select `System` > `Set Clipboard`
1. In the Text box put `%message` and select the back button to save the step.
1. Finally, we need to close the if statement so add a `Task` > `End If` step and then hit the back button to save the Task.
1. Select the checkmark on the Tasker top bar to save your work.

### Sending your clipboard to clipshift

On the latest versions of Android a share menu will appear when the clipboard is set, when it appears just share to ntfy and select the clipshift topic.

There is a *very* advanced way to have Tasker monitor your clipboard but it requires some hacky stuff, I'll type up some instructions if people ask for it.