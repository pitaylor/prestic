#!/usr/bin/env python3

import json


def validate(d, complete=True, **types):
    assert type(d) is dict, "expected dict"
    for k, v in d.items():
        if k in types.keys():
            assert type(v) is types[k], f"'{k}' must be a {types[k].__name__}"
        elif complete and not k.startswith("//"):
            assert False, "'{}' is an unsupported key".format(k)


class Preset:
    def __init__(self, env=None, flags=None, args=None):
        self.env = env or {}
        self.flags = flags or {}
        self.args = args or []

    def __add__(self, o):
        return Preset(
            env={**self.env, **o.env},
            flags={**self.flags, **o.flags},
            args=self.args + o.args,
        )

    def full_args(self):
        result = []
        for k, v in self.flags.items():
            result.extend(self.expand_flag(k, v))
        result.extend(self.args)
        return result

    @staticmethod
    def expand_flag(name, value):
        if not name.startswith("-"):
            name = "--" + name

        if value is True:
            return [name]

        assert type(value) in (int, str), f"invalid flag value: {value}"
        return [name, str(value)]

    @classmethod
    def from_dict(cls, d):
        validate(d, env=dict, flags=dict, args=list)

        return cls(env=d.get("env"), flags=d.get("flags"), args=d.get("args"))


class Config:
    def __init__(self, presets=None, backups=None, forgets=None, prune=None):
        self.presets = presets or {}
        self.backups = backups or []
        self.forgets = forgets or []
        self.prune = prune

    @classmethod
    def from_dict(cls, d):
        validate(d, presets=dict, backups=list, forgets=list, prune=dict)

        presets = {k: Preset.from_dict(v) for k, v in d.get("presets", {}).items()}

        return cls(
            presets=presets,
            backups=[cls.build_preset(presets, v) for v in d.get("backups", [])],
            forgets=[cls.build_preset(presets, v) for v in d.get("forgets", [])],
            prune=cls.build_preset(presets, d["prune"]) if "prune" in d else None,
        )

    @classmethod
    def build_preset(cls, presets, d):
        preset = Preset()

        if "preset" in d:
            validate(d, complete=False, preset=str)
            for name in d.pop("preset").split(","):
                name = name.strip()
                assert name in presets, f"unknown preset: {name}"
                preset += presets[name]

        return preset + Preset.from_dict(d)


class Backups:
    def __init__(self, config):
        self.config = config

    def commands(self):
        for backup in self.config.backups:
            yield backup.env, ["restic", "backup", *backup.full_args()]

        for forget in self.config.forgets:
            yield forget.env, ["restic", "forget", *forget.full_args()]

        if prune := self.config.prune:
            yield prune.env, ["restic", "prune", *prune.full_args()]


raw = json.loads(
    """
{
    "presets": {
        "cloud": {
            "env": {
                "B2_ACCOUNT_ID": "abc",
                "B2_ACCOUNT_KEY": "xyz"
            }
        }
    },
    "//this is a comment": "",
    "backups": [
        {
            "preset": "cloud",
            "args": [
                "~/Documents",
                "~/Pictures"
            ],
            "flags": {
                "stdin": true,
                "exclude": "xyz"
            }
        },
        {
            "args": [
                "~/Documents",
                "~/Pictures"
            ],
            "flags": {
                "stdin": true,
                "exclude": "xyz"
            }
        }
    ],
    
    "prune": {
        "env": {
            "FOO": "BAR"
        }
    }
}
"""
)


b = Backups(Config.from_dict(raw))

for cmd in b.commands():
    print(cmd)
