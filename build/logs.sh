#!/bin/bash

serviceName="<serviceName>"

kubectl logs service/${serviceName} --follow=true