package recorder

import (
	"context"
	"dniprom-cli/internal/container"
	"dniprom-cli/pkg/logger"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Recorder interface {
	PutRich(columns []RichText) error
}

type recorder struct {
	service   *sheets.Service
	container container.Container
	cursor    int64
}

func NewRecorder(ctx context.Context, container container.Container) (Recorder, error) {
	service, err := sheets.NewService(
		ctx,
		option.WithCredentialsFile(container.GetConfig().GoogleCredentials),
	)
	if err != nil {
		return nil, err
	}
	return &recorder{
		service:   service,
		container: container,
		cursor:    0,
	}, nil
}

func (r *recorder) PutRich(columns []RichText) error {
	log := r.container.GetLogger()
	fileID := r.container.GetConfig().FileID
	var cells = make([]*sheets.CellData, 0, len(columns))

	for _, column := range columns {
		bgColor := convertToSpreadsheetColor(column.BackgroundColor)
		cell := sheets.CellData{
			UserEnteredValue: &sheets.ExtendedValue{
				StringValue: &column.Value,
			},
			UserEnteredFormat: &sheets.CellFormat{
				BackgroundColor: bgColor,
			},
			TextFormatRuns: []*sheets.TextFormatRun{
				{
					StartIndex: 0,
					Format: &sheets.TextFormat{
						Bold: column.IsBold,
					},
				},
			},
		}
		if column.Link != "" {
			formulaLink := fmt.Sprintf(`=HYPERLINK("%s","%s")`, column.Link, column.Value)
			cell.UserEnteredValue = &sheets.ExtendedValue{
				FormulaValue: &formulaLink,
			}
			cell.TextFormatRuns = nil
		}
		cells = append(cells, &cell)
	}
	req := sheets.Request{
		UpdateCells: &sheets.UpdateCellsRequest{
			Rows: []*sheets.RowData{
				{
					Values: cells,
				},
			},
			Fields: "*",
			Start: &sheets.GridCoordinate{
				ColumnIndex: 0,
				RowIndex:    r.cursor,
			},
		},
	}
	r.cursor += 1
	_, err := r.service.Spreadsheets.BatchUpdate(fileID, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			&req,
		},
	}).Do()
	if err != nil {
		log.Error(
			"failed to update spreadsheet",
			logger.F("fileID", fileID),
			logger.FError(err),
		)
		return err
	}
	return nil
}

func convertToSpreadsheetColor(color *Color) *sheets.Color {
	if color == nil {
		return nil
	}
	return &sheets.Color{
		Red:   color.Red,
		Green: color.Green,
		Blue:  color.Blue,
	}
}
