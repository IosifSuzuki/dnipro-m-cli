package command

import (
	"dniprom-cli/internal/client"
	"dniprom-cli/internal/container"
	"dniprom-cli/internal/service/recorder"
	"dniprom-cli/internal/worker"
	"dniprom-cli/pkg/logger"
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

type WarrantyCommand struct {
	container    container.Container
	dniproClient client.DniproClient
	recorder     recorder.Recorder
}

func NewWarrantyCommand(container container.Container, client client.DniproClient, recorder recorder.Recorder) *WarrantyCommand {
	return &WarrantyCommand{
		container:    container,
		dniproClient: client,
		recorder:     recorder,
	}

}

func (w *WarrantyCommand) Run(cmd *cobra.Command, args []string) {
	log := w.container.GetLogger()
	config := w.container.GetConfig()
	warrantyWorker := worker.NewWarrantyWorker(w.container, w.dniproClient)
	delay := 500 * time.Millisecond // half a second

	startAt := time.Now().UTC()

	yellowColor := recorder.Color{
		Red:   1,
		Green: 1,
		Blue:  0,
	}

	//Headers
	IDHeaderRichText := recorder.RichText{
		Value:           "ID",
		IsBold:          true,
		BackgroundColor: &yellowColor,
	}
	CodeHeaderRichText := recorder.RichText{
		Value:           "Product Code",
		IsBold:          true,
		BackgroundColor: &yellowColor,
	}
	TitleHeaderRichText := recorder.RichText{
		Value:           "Title",
		IsBold:          true,
		BackgroundColor: &yellowColor,
	}
	WarrantyHeaderRichText := recorder.RichText{
		Value:           "Warranty",
		IsBold:          true,
		BackgroundColor: &yellowColor,
	}
	err := w.recorder.PutRich([]recorder.RichText{
		IDHeaderRichText,
		CodeHeaderRichText,
		TitleHeaderRichText,
		WarrantyHeaderRichText,
	})
	if err != nil {
		log.Error("fail to record header", logger.FError(err))
	}
	for _, productCode := range config.ProductCodes {
		warranty, err := warrantyWorker.FetchByCode(productCode)
		if err != nil {
			log.Error(
				"fail to fetch warranty by code",
				logger.FError(err),
				logger.F("productCode", productCode),
			)
		}
		IDRichText := recorder.RichText{
			Value: fmt.Sprintf("%d", warranty.ID),
		}
		CodeRichText := recorder.RichText{
			Value: productCode,
		}
		TitleRichText := recorder.RichText{
			Value: warranty.Title,
		}
		WarrantyRichText := recorder.RichText{
			Value: warranty.WarrantyText,
		}
		err = w.recorder.PutRich([]recorder.RichText{
			IDRichText,
			CodeRichText,
			TitleRichText,
			WarrantyRichText,
		})
		if err != nil {
			log.Error(
				"fail to record warranty info in row",
				logger.FError(err),
				logger.F("productCode", productCode),
			)
			break
		}
		time.Sleep(delay)
	}
	endAt := time.Now().UTC()
	StartAtTextRichText := recorder.RichText{
		Value: "Start at: ",
	}
	StartAtValueRichText := recorder.RichText{
		Value:  startAt.Format(time.DateTime),
		IsBold: true,
	}
	EndAtTextRichText := recorder.RichText{
		Value: "End at: ",
	}
	EndAtValueRichText := recorder.RichText{
		Value:  endAt.Format(time.DateTime),
		IsBold: true,
	}
	err = w.recorder.PutRich([]recorder.RichText{
		StartAtTextRichText, StartAtValueRichText,
	})
	if err != nil {
		log.Error("fail to record start date info", logger.FError(err))
		return
	}
	time.Sleep(delay)
	err = w.recorder.PutRich([]recorder.RichText{
		EndAtTextRichText, EndAtValueRichText,
	})
	if err != nil {
		log.Error("fail to record end date info", logger.FError(err))
		return
	}
	time.Sleep(delay)

	err = w.recorder.PutRich([]recorder.RichText{
		{
			Value: "Powered by",
		},
		{
			Value: "iOSmates",
			Link:  "https://iosmates.com",
		},
	})
	if err != nil {
		log.Error("fail to record footer info", logger.FError(err))
		return
	}
}
