# shinpuru's Permission System *(since v.0.14)*

## Preamble

Before version 0.14, shinpurus permission system was based on permission levels. Every command had its own, predefined level which a user must had at least to perform this command. Because this system is unflexible and hardly static, I've decided to step over to a domain based permission system.

## Domain Based Permissions

A permission domain looks like following, for example:

```
sp.chat.vote.close
^  ^    ^    ^
|  |    |    |---- sub command
|  |    |--------- main command
|  |-------------- command domain
|----------------- main shinpuru domain
```

So every command and some sub command which need seperate permission configurations are grouped by their command domains and command names.

## Rules

Now, you can specify rules to Discord guild roles which look like following:

```
+sp.guild.config.modlog
^\--------------------/
|          |
|          |----- permission domain
|---------------- rule type specification
```

The `rule type specification` is wether `+` for `ALLOW` or `-` for `DENY`. One of both specifiers **must** be set for each rule.

Of course, you can also specify rules over domain groups by using wild cards:

```
+sp.guild.config.*
```

Rules are always bound to guild roles. So, shinpuru displays role configurations like following:

```
@everyone
    +sp.etc.*
    +sp.chat.*
    -sp.chat.vote.close

@Moderator
    +sp.chat.vote.close
    +sp.guild.mod.*
    -sp.guild.mod.ban

@Admin
    +sp.guild.mod.ban
    +sp.guild.config.*
```

Also keep in mind if a user has an `ALLOW`ing and `DENY`ing rules over the same domain at the same time, this will be counted like the rule is not set at all.

So if a member has following roles with following rules:
```
@Member
    +sp.chat.*
    -sp.chat.vote.close

@Moderator
    +sp.chat.vote.close
```

Their resulting rule set looks like following:
```
+sp.chat.*
```

So this user will be able to use the `sp.chat.vote.close` command because they have access to all commands and sub commands from the domain `sp.chat`.

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
