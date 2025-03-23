package service

import apicontext "github.com/softwareplace/http-utils/context"

type FileService struct {
}

func (s *FileService) UploadFileRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.BadRequest("Failed to upload file")
}
