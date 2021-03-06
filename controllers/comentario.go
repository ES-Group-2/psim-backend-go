package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"poscomp-simulator.com/backend/models"
	"poscomp-simulator.com/backend/utils"
)

func (a *App) GetComentariosSinalizados(w http.ResponseWriter, r *http.Request) {

	if ok, _ := utils.AuthUser(a.DB, w, r, 1); !ok {
		return
	}

	var bc models.BatchComentarios
	if err := bc.GetComentariosSinalizados(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, bc)

}

func (a *App) GetComentariosQuestao(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var err error
	var bc models.BatchComentarios

	if id, ok := vars["id"]; ok {
		bc.QuestaoID, err = strconv.Atoi(id)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
			return
		}
	} else {
		utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
		return
	}

	if err := bc.GetComentariosQuestao(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, bc)

}

func (a *App) PostComentarioQuestao(w http.ResponseWriter, r *http.Request) {

	ok, user := utils.AuthUser(a.DB, w, r, 0)
	if !ok {
		return
	}

	var c models.Comentario
	vars := mux.Vars(r)
	var err error

	c.AutorID = user.Email

	if id, ok := vars["id"]; ok {
		c.QuestaoID, err = strconv.Atoi(id)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
			return
		}
	} else {
		utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := c.Post(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (a *App) ReportComentario(w http.ResponseWriter, r *http.Request) {

	var c models.Comentario
	vars := mux.Vars(r)
	var err error

	if id, ok := vars["id"]; ok {
		c.ID, err = strconv.Atoi(id)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
			return
		}
	} else {
		utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
		return
	}

	if err := c.Report(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (a *App) CleanComentario(w http.ResponseWriter, r *http.Request) {

	ok, _ := utils.AuthUser(a.DB, w, r, 1)
	if !ok {
		return
	}

	var c models.Comentario
	vars := mux.Vars(r)
	var err error

	if id, ok := vars["id"]; ok {
		c.ID, err = strconv.Atoi(id)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
			return
		}
	} else {
		utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
		return
	}

	if err := c.Clean(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (a *App) DeleteComentario(w http.ResponseWriter, r *http.Request) {

	ok, user := utils.AuthUser(a.DB, w, r, 0)
	if !ok {
		return
	}

	var c models.Comentario
	vars := mux.Vars(r)
	var err error

	if id, ok := vars["id"]; ok {
		c.ID, err = strconv.Atoi(id)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
			return
		}
	} else {
		utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
		return
	}

	if err := c.Get(a.DB); err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Coment??rio n??o foi encontrado.")
			return
		}
		utils.RespondWithError(w, http.StatusBadRequest, "N??o foi poss??vel apagar o coment??rio.")
		return
	}

	if user.Email == c.AutorID || user.NivelAcesso > 0 {
		c.Delete(a.DB)
		w.WriteHeader(http.StatusOK)
		return
	}

	utils.RespondWithError(w, http.StatusUnauthorized, "Coment??rio n??o pertence ao usu??rio ou n??vel de acesso insuficiente.")

}
