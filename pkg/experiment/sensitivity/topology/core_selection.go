// Copyright (c) 2017 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package topology

import (
	"github.com/intelsdi-x/swan/pkg/isolation/topo"
	"github.com/intelsdi-x/swan/pkg/utils/errutil"
)

// sharedCacheThreads returns threads from one socket that share a last-level
// cache. To avoid placing workloads on both hyperthreads for any physical
// core, only one thread from each is included in the result.
func sharedCacheThreads() topo.ThreadSet {
	allThreads, err := topo.Discover()
	errutil.Check(err)

	// Retain only threads for one socket.
	socket, err := allThreads.Sockets(1)
	errutil.Check(err)

	// Retain only one thread per physical core.
	// NB: The following filter prediccate closes over this int set.
	temp := socket.AvailableCores()

	return socket.Filter(func(t topo.Thread) bool {
		retain := temp.Contains(t.Core())
		temp.Remove(t.Core())
		return retain
	})
}

func getSiblingThreadsOfThread(reservedThread topo.Thread) topo.ThreadSet {
	requestedCore := reservedThread.Core()

	allThreads, err := topo.Discover()
	errutil.Check(err)

	threadsFromCore, err := allThreads.FromCores(requestedCore)
	errutil.Check(err)

	return threadsFromCore.Remove(reservedThread)
}

func getSiblingThreadsOfThreadSet(threads topo.ThreadSet) (results topo.ThreadSet) {
	for _, thread := range threads {
		siblings := getSiblingThreadsOfThread(thread)
		for _, sibling := range siblings {
			// Omit the reserved threads; if at least one pair from threads were
			// siblings of each other, they would both otherwise wrongly end up in
			// the result.
			if !threads.Contains(sibling) {
				results = append(results, sibling)
			}
		}
	}

	return results
}
