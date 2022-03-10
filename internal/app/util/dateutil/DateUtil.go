// パッケージを個別のものにし、疑似的にstaticのような運用。
package dateutil

import (
	"log"
	"time"
)

// Golangのフォーマット文字列は、他の言語と違うので注意。
// @see https://go.dev/src/time/format.go

// 現在日時をstringで返す。
func GetNowString() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

// string型の日付をtime.Timeに変換する。
// 第二戻り値に成否情報をもたせている。
func ParseStringToTime(datetime string) (time.Time, bool) {
	parseDatetime, err := time.Parse("2006/01/02 15:04:05", datetime)
	if err != nil {
		log.Printf("Error ParseStringToTime. value: %s", datetime)
		return parseDatetime, false
	}

	return parseDatetime, true
}

// beforeとafterの差を秒数で返す。
// beforeもafterも "2006/01/02 15:04:05" の形式であることが前提。
// 第二戻り値に成否情報をもたせている。
func DiffSecond(before string, after string) (float64, bool) {
	convBefore, ok := ParseStringToTime(before)
	if !ok {
		log.Printf("Error DiffSecond. value: %s", before)
		return 0.0, false
	}
	convAfter, ok := ParseStringToTime(after)
	if !ok {
		log.Printf("Error DiffSecond. value: %s", after)
		return 0.0, false
	}

	return convAfter.Sub(convBefore).Seconds(), true
}
