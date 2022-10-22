package queries

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math"
	"time"

	"github.com/AgileProggers/archiv-backend-go/pkg/database"
	"github.com/AgileProggers/archiv-backend-go/pkg/models"
	"gorm.io/gorm"
)

func GetAllVods(v *[]models.Vod, query models.Vod, pagination Pagination, o string) (*Pagination, error) {
	if o == "" {
		o = "date desc"
	}
	result := database.DB.Model(&query).Omit("transcript")
	if query.Title != "" {
		// if title is given, do case insensitive search in title string
		result = result.Where("position(LOWER(?) in LOWER(title))>0", query.Title)
	} else {
		// else search exact query match
		result = result.Where(query)
	}
	result = result.Where("publish = ?", true).Order(o).Find(v).Scopes(Paginate(len(*v), &pagination, database.DB)).Find(v)
	if result.RowsAffected == 0 {
		return &pagination, errors.New("not found")
	}
	return &pagination, nil
}

func AddNewVod(v *models.Vod) error {
	if v.UUID == "" {
		b := make([]byte, 3)
		for {
			rand.Read(b)
			uuid := hex.EncodeToString(b)
			if err := GetOneVod(&models.Vod{}, uuid, false); err != nil {
				v.UUID = uuid
				break
			}
		}
	}

	if err := database.DB.Create(v).Error; err != nil {
		return err
	}

	// update vod timestamp
	var settings models.Settings
	settings.DateVodsUpdate = time.Now()
	if err := PartiallyUpdateSettings(&settings); err != nil {
		return err
	}

	return nil
}

func GetOneVod(v *models.Vod, uuid string, onlyPublic bool) error {
	var result *gorm.DB
	if onlyPublic {
		result = database.DB.Where("uuid = ?", uuid).Where("publish = ?", true).Preload("Clips.Creator").Preload("Clips.Game")
	} else {
		result = database.DB.Where("uuid = ?", uuid)
	}
	result = result.Find(v)
	if result.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

func GetVodByFilename(v *models.Vod, filename string) error {
	result := database.DB.Omit("transcript").Where("filename = ?", filename).Find(v)
	if result.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

func GetVodsByUUID(v *[]models.Vod, uuids []string) error {
	result := database.DB.Omit("transcript").Where("uuid IN ?", uuids).Where("publish = ?", true).Find(v)
	if result.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

func GetVodsByYear(v *[]models.Vod, year string) error {
	result := database.DB.Model(&v).Omit("transcript").Where("date_part('year', date) = ?", year).Where("publish = ?", true).Order("date desc").Find(v)
	if result.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

func PatchVod(changes map[string]interface{}, uuid string) error {
	var vod models.Vod
	if err := GetOneVod(&vod, uuid, false); err != nil {
		return errors.New("vod not found")
	}

	if err := database.DB.Model(&vod).Where("uuid = ?", uuid).Updates(changes).Error; err != nil {
		return errors.New("update failed")
	}

	// update vod timestamp
	var settings models.Settings
	settings.DateVodsUpdate = time.Now()
	if err := PartiallyUpdateSettings(&settings); err != nil {
		return err
	}

	return nil
}

func DeleteVod(v *models.Vod, uuid string) error {
	database.DB.Where("uuid = ?", uuid).Delete(v)
	return nil
}

func GetVodsFullText(foundVods *[]map[string]interface{}, query string, pagination Pagination) (*Pagination, error) {
	// ts_headline() is incredibly slow on large text, so we use 2 queries to limit the rows
	var rowCount int64
	database.DB.Model(&models.Vod{}).Select("vods.uuid").Where("publish = true and vods.title_vector @@ websearch_to_tsquery('german', ?) or vods.title_vector @@ websearch_to_tsquery('english', ?) or vods.transcript_vector @@ websearch_to_tsquery('german', ?) or vods.transcript_vector @@ websearch_to_tsquery('english', ?)", query, query, query, query).Count(&rowCount)

	if rowCount == 0 {
		return &pagination, errors.New("not found")
	}

	pagination.TotalRows = rowCount
	pagination.TotalPages = int(math.Ceil(float64(rowCount) / float64(pagination.GetLimit())))

	// formated sql query for GetVodsFullText

	// with chunk as (
	// 	select vods.uuid, vods.title, vods.filename, vods.resolution, vods.fps, vods.size, vods.date, vods.transcript, vods.title_vector, vods.transcript_vector,
	// 		coalesce(vods.duration, 0) as duration,
	// 		coalesce(vods.viewcount, 0) as viewcount,
	// 		coalesce(ts_rank(vods.title_vector, german) + ts_rank(vods.title_vector, english), 0) as title_rank,
	// 		coalesce(ts_rank(vods.transcript_vector, german) + ts_rank(vods.transcript_vector, english), 0) as transcript_rank
	// 	FROM vods,
	// 	websearch_to_tsquery('german', 'arch') as german,
	// 	websearch_to_tsquery('english', 'arch') as english
	// 	WHERE publish = true and vods.title_vector @@ german or vods.title_vector @@ english or vods.transcript_vector @@ german or vods.transcript_vector @@ english
	// 	ORDER BY title_rank desc, transcript_rank desc
	// )
	// select chunk.uuid, chunk.title, chunk.filename, chunk.duration, chunk.viewcount, chunk.resolution, chunk.fps, chunk.size, chunk.date, chunk.title_rank, chunk.transcript_rank,
	// coalesce(ts_headline(chunk.title, german && english, 'MaxFragments=6, StartSel=<span>, StopSel=</span>, FragmentDelimiter=<hr>'), '') as title_matches,
	// coalesce(ts_headline(chunk.transcript, german && english, 'MaxFragments=6, StartSel=<span>, StopSel=</span>, FragmentDelimiter=<hr>'), '') as transcript_matches
	// from chunk, websearch_to_tsquery('german', 'arch') as german, websearch_to_tsquery('english', 'arch') as english
	// where chunk.title_vector @@ german or chunk.title_vector @@ english or chunk.transcript_vector @@ german or chunk.transcript_vector @@ english
	// limit 4 offset 8

	result := database.DB.Raw("with chunk as (select vods.uuid, vods.title, vods.filename, vods.resolution, vods.fps, vods.size, vods.date, vods.transcript, vods.title_vector, vods.transcript_vector, coalesce(vods.duration, 0) as duration, coalesce(vods.viewcount, 0) as viewcount, coalesce(ts_rank(vods.title_vector, german) + ts_rank(vods.title_vector, english), 0) as title_rank, coalesce(ts_rank(vods.transcript_vector, german) + ts_rank(vods.transcript_vector, english), 0) as transcript_rank FROM vods, websearch_to_tsquery('german', ?) as german, websearch_to_tsquery('english', ?) as english WHERE publish = true and vods.title_vector @@ german or vods.title_vector @@ english or vods.transcript_vector @@ german or vods.transcript_vector @@ english ORDER BY title_rank desc, transcript_rank desc) select chunk.uuid, chunk.title, chunk.filename, chunk.duration, chunk.viewcount, chunk.resolution, chunk.fps, chunk.size, chunk.date, chunk.title_rank, chunk.transcript_rank, coalesce(ts_headline(chunk.title, german && english, 'MaxFragments=6, StartSel=<span>, StopSel=</span>, FragmentDelimiter=<hr>'), '') as title_matches, coalesce(ts_headline(chunk.transcript, german && english, 'MaxFragments=6, StartSel=<span>, StopSel=</span>, FragmentDelimiter=<hr>'), '') as transcript_matches	from chunk, websearch_to_tsquery('german', ?) as german, websearch_to_tsquery('english', ?) as english where chunk.title_vector @@ german or chunk.title_vector @@ english or chunk.transcript_vector @@ german or chunk.transcript_vector @@ english limit ? offset ?", query, query, query, query, pagination.Limit, pagination.GetOffset()).Find(foundVods)

	if result.RowsAffected == 0 {
		return &pagination, errors.New("not found")
	}

	return &pagination, nil
}
