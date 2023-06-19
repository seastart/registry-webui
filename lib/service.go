package lib

import (
	"errors"
	"sort"
	"strings"
	"sync"
)

// repoes service
var Repoes = reposervice{}

type reposervice struct {
	sync.RWMutex
	repoes []*Repo // repoes
}

// refresh all repoes
func (r *reposervice) Refresh() error {
	logger.Warnf("refresh all repoes")
	r.Lock()
	defer r.Unlock()
	repoes, err := getAllRepos()
	if err != nil {
		logger.Errorf("refresh all repoes got error: %s", err.Error())
		return err
	} else {
		r.repoes = repoes
		r.resort()
	}
	return nil
}

// get page date
func (r *reposervice) GetPage(keyword string, page int, pageSize int) (repoes []*Repo, hasmore bool, err error) {
	r.RLock()
	defer r.RUnlock()
	// if never fetched, return error
	if r.repoes == nil {
		return nil, false, errors.New("repoes not fetched, please wait system init")
	}
	if keyword != "" {
		for _, repo := range r.repoes {
			if strings.Contains(repo.Name, keyword) || strings.Contains(repo.Desc, keyword) {
				repoes = append(repoes, repo)
			}
		}
	} else {
		repoes = r.repoes
	}
	// subarray
	start := (page - 1) * pageSize
	end := page * pageSize
	if start < 0 || start >= len(repoes) {
		return nil, false, nil
	} else if end > len(repoes) {
		end = len(repoes)
	}
	if start == len(repoes) {
		return nil, false, nil
	} else if end == len(repoes) {
		return repoes[start:end], false, nil
	} else {
		return repoes[start:end], true, nil
	}
}

// get repo detail
func (r *reposervice) GetDetail(name string, refresh bool) (*Repo, error) {
	if refresh {
		repo, err := getRepoDetail(name)
		if err != nil {
			return nil, err
		}
		// update to list
		r.Lock()
		has := false
		for i, repo := range r.repoes {
			if repo.Name == name {
				r.repoes[i] = repo
				has = true
				break
			}
		}
		if !has {
			r.repoes = append(r.repoes, repo)
		}
		r.resort()
		r.Unlock()
		return repo, nil
	} else {
		r.RLock()
		for _, repo := range r.repoes {
			if repo.Name == name {
				r.RUnlock()
				return repo, nil
			}
		}
		r.RUnlock()
		// try fresh if not exsits
		return r.GetDetail(name, true)
	}
}

// resort by updated_at
func (r *reposervice) resort() {
	sort.Slice(r.repoes, func(i, j int) bool {
		return r.repoes[i].LastUpdate > r.repoes[j].LastUpdate
	})
}
