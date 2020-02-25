# shinpuru's Permission System *(since v.0.14)*

## Preamble

Before version 0.14, shinpuru's permission system was based on permission levels. Every command had a specific, predefined and static level which a user must match at least to perform the command. Because of this system's unflexibility, I've decided to step over to a domain based permission system.

## Domain Based Permissions

A permission domain looks like following, for example:

```
sp.chat.vote.close
^  ^    ^    ^
|  |    |    +---- sub command
|  |    +--------- main command
|  +-------------- command domain
+----------------- main shinpuru domain
```

So every command and some sub command which need seperate permission configurations are grouped by their command domains and command names.

## Rules

Now, you can specify rules to Discord guild roles which look like following:

```
+sp.guild.config.modlog
^\--------------------/
|          |
|          +----- permission domain
+---------------- rule type specification
```

The `rule type specification` is either `+` for `ALLOW` or `-` for `DENY`. One of both specifiers **must** be set for each rule. `DENY` always counts over `ALLOW` on the same permission level.

Of course, you can also specify rules over domain groups by using wild cards:

```
+sp.guild.config.*
```

Keep in mind that higer domains count over lower when using wildcards, which means `-sp.guild.mod.kick` counts over `+sp.guild.mod.*` and denies the usage of the `kick` command.

Rules are always bound to guild roles. So, shinpuru displays role configurations like following:

```
@Admin
    +sp.guild.mod.ban
    +sp.guild.config.*

@Moderator
    +sp.chat.vote.close
    +sp.guild.mod.*
    -sp.guild.mod.ban

@everyone
    +sp.etc.*
    +sp.chat.*
```

If you want to negate a rule, just set the negative rule on the specific role.

So if you have a role with following rule set:

```
@Moderator
    +sp.chat.vote.close
    +sp.guild.mod.*
    -sp.guild.mod.ban
```

And now, you want to allow `@Moderators` to use the `ban` command, just set the rule

```
+sp.guild.mod.ban
```

which results in the follwoing rule set configuration:

```
@Moderator
    +sp.chat.vote.close
    +sp.guild.mod.*
```

## Rule Domination Order

Generally, higer group levels in the domain override rules with lower levels.

Lets take a look at the following example rule set:

```
@Moderator
    +sp.guild.mod.*
    -sp.guild.mod.ban
```

The user assigned to this role will be able to use all commands of `sp.guild.mod` with exception of the `sp.guild.mod.ban` command which will be denied.

This also works vice versa:

```
@Moderator
    -sp.guild.config.*
    +sp.guild.config.autorole
```

This user will not be able to use any command from `sp.guild.config` except the `sp.guild.config.autorole` command.

Also, rules et to roles with higer position will always override rules assigned to troles with lower position, like in the following example:

*Role position equals like displayed. Higher means higher role position.*
```
@Supporter
    -sp.chat.vote.close

@Moderator
    +sp.chat.vote.close
    +sp.guild.mod.*
```

If a user has both roles, `@Supporter` and `@Moderator`, the rule `+sp.chat.vote.close` of `@Moderator` will be canceled out by the rule `-sp.chat.vote.close` bound to `@Supporter`, which is higer in position.