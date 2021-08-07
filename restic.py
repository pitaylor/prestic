#!/usr/bin/env python3

import json
import sys


def validate(obj, inclusive=True, **types):
    assert type(obj) is dict, "expected dict"
    for k, v in obj.items():
        if k in types.keys():
            assert type(v) is types[k], f"'{k}' must be a {types[k].__name__}"
        elif inclusive and not k.startswith("//"):
            assert False, "'{}' is an invalid key".format(k)


class Preset:
    def __init__(self, env=None, flags=None, args=None):
        self.env = env or {}
        self.flags = flags or {}
        self.args = args or []

    def __add__(self, other):
        return Preset(
            env={**self.env, **other.env},
            flags={**self.flags, **other.flags},
            args=self.args + other.args,
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
    def from_dict(cls, obj):
        validate(obj, env=dict, flags=dict, args=list)

        return cls(env=obj.get("env"), flags=obj.get("flags"), args=obj.get("args"))


class Config:
    def __init__(self, presets=None, backups=None, forgets=None, prune=None):
        self.presets = presets or {}
        self.backups = backups or []
        self.forgets = forgets or []
        self.prune = prune

    @classmethod
    def from_dict(cls, obj):
        validate(obj, presets=dict, backups=list, forgets=list, prune=dict)

        presets = {k: Preset.from_dict(v) for k, v in obj.get("presets", {}).items()}

        return cls(
            presets=presets,
            backups=[cls.build_preset(presets, v) for v in obj.get("backups", [])],
            forgets=[cls.build_preset(presets, v) for v in obj.get("forgets", [])],
            prune=cls.build_preset(presets, obj["prune"]) if "prune" in obj else None,
        )

    @staticmethod
    def build_preset(presets, obj):
        preset = Preset()

        if "preset" in obj:
            validate(obj, inclusive=False, preset=str)
            for name in obj.pop("preset").split(","):
                name = name.strip()
                assert name in presets, f"unknown preset: {name}"
                preset += presets[name]

        return preset + Preset.from_dict(obj)


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


config_json = """
{
    "presets": {
        "cloud": {
            "env": {
                "B2_ACCOUNT_ID": "abc",
                "B2_ACCOUNT_KEY": "xyz"
            }
        }
    },
    
    "// this is a comment": "",
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

try:
    config = Config.from_dict(json.loads(config_json))
    backup = Backups(config)
except (AssertionError, json.decoder.JSONDecodeError) as e:
    print(f"ERROR: {e}", file=sys.stderr)
    sys.exit(1)

for cmd in backup.commands():
    print(cmd)
