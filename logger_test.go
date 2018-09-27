package log_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strconv"
	"github.com/pawel0987/fakes"
	"github.com/pawel0987/utils/json_utils"
	"github.com/pawel0987/log"
)

var _ = Describe("Log Session", func() {
	var session log.Logger
	var writer *fakes.FakeWriter

	const (
		GLOBAL_SESSION_NAME = "LoggerTest"
		SUB_SESSION_NAME = "subSession"

		TEST_STRING = "hello"
		TEST_NUMBER = 1
	)
	var (
		TEST_STRING_ARRAY = [...]string{"a", "b"}
		TEST_OBJ = struct {
			Name string `json:"name"`
			Age  int     `json:"age"`
		}{
			"ziom",
			21,
		}
	)

	BeforeEach(func() {
		writer = fakes.NewFakeWriter()
		log.GetConfig().SetOutputSinks(writer, GinkgoWriter)

		session = log.Session(GLOBAL_SESSION_NAME)
	})



	Describe("Creating sessions", func() {
		Context("When format is JSON", func() {
			It("Should log in current session", func() {
				session.Info("test")
				logStream := string(writer.GetContent())
				Expect(logStream).To(ContainSubstring("\"session\":\"" + GLOBAL_SESSION_NAME + "\""))
			})

			Context("When sub-session opened", func() {
				var subSession log.Logger
				BeforeEach(func() {
					subSession = log.Session(SUB_SESSION_NAME)
				})

				It("Should log in sub-session", func() {
					subSession.Info("test")
					logStream := string(writer.GetContent())
					Expect(logStream).To(ContainSubstring("\"session\":\"" + SUB_SESSION_NAME + "\""))
				})
			})
		})
	})

	Describe("Logging", func() {
		Context("When log level is fatal", func() {
			BeforeEach(func() {
				log.GetConfig().SetLogLevel(log.LEVEL_FATAL)
			})

			It("Should log only fatal logs", func() {
				session.Fatal("some message")
				session.Error("some message")
				session.Warning("some message")
				session.Info("some message")
				session.Debug("some message")

				logStream := string(writer.GetContent())
				Expect(logStream).To(ContainSubstring(log.LEVEL_FATAL.String()))
				Expect(logStream).NotTo(ContainSubstring(log.LEVEL_ERROR.String()))
				Expect(logStream).NotTo(ContainSubstring(log.LEVEL_WARNING.String()))
				Expect(logStream).NotTo(ContainSubstring(log.LEVEL_INFO.String()))
				Expect(logStream).NotTo(ContainSubstring(log.LEVEL_DEBUG.String()))
			})
		})

		Context("When log level is error", func() {
			BeforeEach(func() {
				log.GetConfig().SetLogLevel(log.LEVEL_ERROR)
			})

			It("Should log only fatal and error logs", func() {
				session.Fatal("some message")
				session.Error("some message")
				session.Warning("some message")
				session.Info("some message")
				session.Debug("some message")

				logStream := string(writer.GetContent())
				Expect(logStream).To(ContainSubstring(log.LEVEL_FATAL.String()))
				Expect(logStream).To(ContainSubstring(log.LEVEL_ERROR.String()))
				Expect(logStream).NotTo(ContainSubstring(log.LEVEL_WARNING.String()))
				Expect(logStream).NotTo(ContainSubstring(log.LEVEL_INFO.String()))
				Expect(logStream).NotTo(ContainSubstring(log.LEVEL_DEBUG.String()))
			})
		})

		Context("When log level is warning", func() {
			BeforeEach(func() {
				log.GetConfig().SetLogLevel(log.LEVEL_WARNING)
			})

			It("Should log only fatal, error and warning logs", func() {
				session.Fatal("some message")
				session.Error("some message")
				session.Warning("some message")
				session.Info("some message")
				session.Debug("some message")

				logStream := string(writer.GetContent())
				Expect(logStream).To(ContainSubstring(log.LEVEL_FATAL.String()))
				Expect(logStream).To(ContainSubstring(log.LEVEL_ERROR.String()))
				Expect(logStream).To(ContainSubstring(log.LEVEL_WARNING.String()))
				Expect(logStream).NotTo(ContainSubstring(log.LEVEL_INFO.String()))
				Expect(logStream).NotTo(ContainSubstring(log.LEVEL_DEBUG.String()))
			})
		})

		Context("When log level is info", func() {
			BeforeEach(func() {
				log.GetConfig().SetLogLevel(log.LEVEL_INFO)
			})

			It("Should log only fatal, error, warning and info logs", func() {
				session.Fatal("some message")
				session.Error("some message")
				session.Warning("some message")
				session.Info("some message")
				session.Debug("some message")

				logStream := string(writer.GetContent())
				Expect(logStream).To(ContainSubstring(log.LEVEL_FATAL.String()))
				Expect(logStream).To(ContainSubstring(log.LEVEL_ERROR.String()))
				Expect(logStream).To(ContainSubstring(log.LEVEL_WARNING.String()))
				Expect(logStream).To(ContainSubstring(log.LEVEL_INFO.String()))
				Expect(logStream).NotTo(ContainSubstring(log.LEVEL_DEBUG.String()))
			})
		})

		Context("When log level is debug", func() {
			BeforeEach(func() {
				log.GetConfig().SetLogLevel(log.LEVEL_DEBUG)
			})

			It("Should log everything", func() {
				session.Fatal("some message")
				session.Error("some message")
				session.Warning("some message")
				session.Info("some message")
				session.Debug("some message")

				logStream := string(writer.GetContent())
				Expect(logStream).To(ContainSubstring(log.LEVEL_FATAL.String()))
				Expect(logStream).To(ContainSubstring(log.LEVEL_ERROR.String()))
				Expect(logStream).To(ContainSubstring(log.LEVEL_WARNING.String()))
				Expect(logStream).To(ContainSubstring(log.LEVEL_INFO.String()))
				Expect(logStream).To(ContainSubstring(log.LEVEL_DEBUG.String()))
			})
		})
	})

	Describe("Logging with data", func() {
		It("Should omit data when empty", func() {
			session.Info("Test", log.Data{})
			logStream := string(writer.GetContent())
			Expect(logStream).NotTo(ContainSubstring("\"data\":["))
		})

		It("Should log with string as data", func() {
			session.Info("Test", log.Data{"something": TEST_STRING})
			logStream := string(writer.GetContent())
			Expect(logStream).To(ContainSubstring("\"data\":{\"something\":\"" + TEST_STRING + "\"}"))
		})

		It("Should log with int as data", func() {
			session.Info("Test", log.Data{"something": TEST_NUMBER})
			logStream := string(writer.GetContent())
			Expect(logStream).To(ContainSubstring("\"data\":{\"something\":" + strconv.Itoa(TEST_NUMBER) + "}"))
		})

		It("Should log with []string as data", func() {
			session.Info("Test", log.Data{"something": TEST_STRING_ARRAY})
			logStream := string(writer.GetContent())
			Expect(logStream).To(ContainSubstring("\"data\":{\"something\":" + json_utils.Encode(TEST_STRING_ARRAY) + "}"))
		})

		It("Should log with interface{} as data", func() {

			session.Info("Test", log.Data{"something": TEST_OBJ})
			logStream := string(writer.GetContent())
			Expect(logStream).To(ContainSubstring("\"data\":{\"something\":" + json_utils.Encode(TEST_OBJ) + "}"))
		})
	})

	Describe("Logging with custom fields", func() {
		It("Should log with string custom field", func() {
			log.GetConfig().SetCustomFields(log.Data{"something": TEST_STRING})

			session.Info("Test")
			logStream := string(writer.GetContent())
			Expect(logStream).To(ContainSubstring("\"something\":\"" + TEST_STRING + "\""))
		})

		It("Should log with int custom field", func() {
			log.GetConfig().SetCustomFields(log.Data{"something": TEST_NUMBER})

			session.Info("Test")
			logStream := string(writer.GetContent())
			Expect(logStream).To(ContainSubstring("\"something\":" + strconv.Itoa(TEST_NUMBER)))
		})

		It("Should log with []string custom field", func() {
			log.GetConfig().SetCustomFields(log.Data{"something": TEST_STRING_ARRAY})

			session.Info("Test")
			logStream := string(writer.GetContent())
			Expect(logStream).To(ContainSubstring("\"something\":" + json_utils.Encode(TEST_STRING_ARRAY)))
		})

		It("Should log with interface{} custom field", func() {
			log.GetConfig().SetCustomFields(log.Data{"something": TEST_OBJ})

			session.Info("Test")
			logStream := string(writer.GetContent())
			Expect(logStream).To(ContainSubstring("\"something\":" + json_utils.Encode(TEST_OBJ)))
		})
	})
})
