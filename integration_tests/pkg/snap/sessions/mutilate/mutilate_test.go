package mutilate

import (
	"path"
	"testing"

	"github.com/intelsdi-x/athena/pkg/snap/sessions/mutilate"
	"github.com/intelsdi-x/athena/pkg/utils/fs"
	"github.com/intelsdi-x/swan/integration_tests/pkg/snap/sessions"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSnapMutilateSession(t *testing.T) {

	Convey("When testing MutilateSnapSession ", t, func() {
		Convey("We have snapd running ", func() {

			cleanupSnap, loader, snapteldAddress := sessions.RunAndTestSnaptel()
			defer cleanupSnap()

			Convey("And we loaded publisher plugin", func() {

				cleanupMerticsFile, publisher, publisherDataFilePath := sessions.PreparePublisher(loader)
				defer cleanupMerticsFile()

				Convey("Then we prepared and launch mutilate session", func() {

					mutilateSessionConfig := mutilatesession.DefaultConfig()
					mutilateSessionConfig.SnapteldAddress = snapteldAddress
					mutilateSessionConfig.Publisher = publisher
					mutilateSnapSession, err := mutilatesession.NewSessionLauncher(mutilateSessionConfig)
					So(err, ShouldBeNil)

					cleanupMockedFile, mockedTaskInfo := sessions.PrepareMockedTaskInfo(path.Join(
						fs.GetSwanPath(), "misc/snap-plugin-collector-mutilate/mutilate/mutilate.stdout"))
					defer cleanupMockedFile()

					handle, err := mutilateSnapSession.LaunchSession(mockedTaskInfo, "foo:bar")
					So(err, ShouldBeNil)

					defer func() {
						err := handle.Stop()
						So(err, ShouldBeNil)
					}()

					Convey("Later we checked if task is running", func() {
						So(handle.IsRunning(), ShouldBeTrue)

						// These are results from test output file
						// in "src/github.com/intelsdi-x/swan/misc/
						// snap-plugin-collector-mutilate/mutilate/mutilate.stdout"
						expectedMetrics := map[string]string{
							"avg":    "20.80000",
							"std":    "23.10000",
							"min":    "11.90000",
							"5th":    "13.30000",
							"10th":   "13.40000",
							"90th":   "33.40000",
							"95th":   "43.10000",
							"99th":   "59.50000",
							"qps":    "4993.10000",
							"custom": "1777.88781",
						}

						Convey("In order to read and test published data", func() {

							dataValid := sessions.ReadAndTestPublisherData(publisherDataFilePath, expectedMetrics)
							So(dataValid, ShouldBeTrue)
						})
					})
				})
			})
		})
	})
}