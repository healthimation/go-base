#!/bin/bash

serviceName="goal"

kubectl logs service/${serviceName} --follow=true
