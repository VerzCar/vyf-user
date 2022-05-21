#!/usr/bin/env bash
#
# Copyright (c) 2019 - 2021. Vecomentman, Carlo Verzeri, all rights reserved.
#

set -e
set -x

mjml ./templates/company/verification/*.mjml --config.minify=y -o ./dist/
mjml ./templates/user/activation/*.mjml --config.minify=y -o ./dist/
mjml ./templates/user/passwordResetDone/*.mjml --config.minify=y -o ./dist/
mjml ./templates/user/resetPassword/*.mjml --config.minify=y -o ./dist/
