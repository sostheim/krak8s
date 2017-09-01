#
# Copyright 2017 Samsung SDSA CNCT
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License. 
#
# Dockerfile - krak8s API service container.
#

FROM quay.io/samsung_cnct/k2:af557b1b081824c0a1121c0c688a3e2085362301
MAINTAINER Rick Sostheim
LABEL vendor="Samsung CNCT"

COPY build/linux_amd64/krak8s /
COPY commands/node_pool.tmpl commands/services.tmpl /
COPY swagger /

WORKDIR "/"
ENTRYPOINT ["/krak8s"]
