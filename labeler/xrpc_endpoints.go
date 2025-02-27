package labeler

import (
	"net/http/httputil"
	"strconv"

	atproto "github.com/bluesky-social/indigo/api/atproto"
	label "github.com/bluesky-social/indigo/api/label"
	"github.com/bluesky-social/indigo/util/version"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
)

func (s *Server) RegisterHandlersComAtproto(e *echo.Echo) error {
	// handle/account hosting
	e.GET("/xrpc/com.atproto.server.describeServer", s.HandleComAtprotoServerDescribeServer)
	// TODO: session create/refresh/delete?

	// minimal moderation reporting/actioning
	e.GET("/xrpc/com.atproto.admin.getModerationAction", s.HandleComAtprotoAdminGetModerationAction)
	e.GET("/xrpc/com.atproto.admin.getModerationActions", s.HandleComAtprotoAdminGetModerationActions)
	e.GET("/xrpc/com.atproto.admin.getModerationReport", s.HandleComAtprotoAdminGetModerationReport)
	e.GET("/xrpc/com.atproto.admin.getModerationReports", s.HandleComAtprotoAdminGetModerationReports)
	e.POST("/xrpc/com.atproto.admin.resolveModerationReports", s.HandleComAtprotoAdminResolveModerationReports)
	e.POST("/xrpc/com.atproto.admin.reverseModerationAction", s.HandleComAtprotoAdminReverseModerationAction)
	e.POST("/xrpc/com.atproto.admin.takeModerationAction", s.HandleComAtprotoAdminTakeModerationAction)
	e.POST("/xrpc/com.atproto.report.create", s.HandleComAtprotoReportCreate)

	// label-specific
	e.GET("/xrpc/com.atproto.label.queryLabels", s.HandleComAtprotoLabelQueryLabels)

	return nil
}

func (s *Server) rewriteProxyRequestAdmin(r *httputil.ProxyRequest) {
	r.SetXForwarded()
	r.SetURL(s.xrpcProxyURL)
	r.Out.Header.Set("Authorization", s.xrpcProxyAuthHeader)
}

func (s *Server) RegisterProxyHandlers(e *echo.Echo) error {

	rp := &httputil.ReverseProxy{Rewrite: s.rewriteProxyRequestAdmin}

	// Proxy some admin requests
	e.GET("/xrpc/com.atproto.admin.getRecord", echo.WrapHandler(rp))
	e.GET("/xrpc/com.atproto.admin.getRepo", echo.WrapHandler(rp))
	e.GET("/xrpc/com.atproto.admin.searchRepos", echo.WrapHandler(rp))

	return nil
}

func (s *Server) HandleComAtprotoServerDescribeServer(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleComAtprotoServerDescribeServer")
	defer span.End()
	var out *atproto.ServerDescribeServer_Output
	var handleErr error
	// func (s *Server) handleComAtprotoServerDescribeServer(ctx context.Context) (*atproto.ServerDescribeServer_Output, error)
	out, handleErr = s.handleComAtprotoServerDescribeServer(ctx)
	if handleErr != nil {
		return handleErr
	}
	return c.JSON(200, out)
}

func (s *Server) HandleComAtprotoLabelQueryLabels(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleComAtprotoLabelQueryLabels")
	defer span.End()
	cursor := c.QueryParam("cursor")

	var limit int
	if p := c.QueryParam("limit"); p != "" {
		var err error
		limit, err = strconv.Atoi(p)
		if err != nil {
			return err
		}
	} else {
		limit = 20
	}

	sources := c.QueryParams()["sources"]

	uriPatterns := c.QueryParams()["uriPatterns"]
	var out *label.QueryLabels_Output
	var handleErr error
	// func (s *Server) handleComAtprotoLabelQueryLabels(ctx context.Context,cursor string,limit int,sources []string,uriPatterns []string) (*comatprototypes.LabelQueryLabels_Output, error)
	out, handleErr = s.handleComAtprotoLabelQueryLabels(ctx, cursor, limit, sources, uriPatterns)
	if handleErr != nil {
		return handleErr
	}
	return c.JSON(200, out)
}

func (s *Server) HandleComAtprotoAdminGetModerationAction(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleComAtprotoAdminGetModerationAction")
	defer span.End()

	id, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		return err
	}
	var out *atproto.AdminDefs_ActionViewDetail
	var handleErr error
	// func (s *Server) handleComAtprotoAdminGetModerationAction(ctx context.Context,id int) (*atproto.AdminDefs_ActionViewDetail, error)
	out, handleErr = s.handleComAtprotoAdminGetModerationAction(ctx, id)
	if handleErr != nil {
		return handleErr
	}
	return c.JSON(200, out)
}

func (s *Server) HandleComAtprotoAdminGetModerationActions(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleComAtprotoAdminGetModerationActions")
	defer span.End()
	before := c.QueryParam("before")

	var limit int
	if p := c.QueryParam("limit"); p != "" {
		var err error
		limit, err = strconv.Atoi(p)
		if err != nil {
			return err
		}
	} else {
		limit = 50
	}
	subject := c.QueryParam("subject")
	var out *atproto.AdminGetModerationActions_Output
	var handleErr error
	// func (s *Server) handleComAtprotoAdminGetModerationActions(ctx context.Context,before string,limit int,subject string) (*atproto.AdminGetModerationActions_Output, error)
	out, handleErr = s.handleComAtprotoAdminGetModerationActions(ctx, before, limit, subject)
	if handleErr != nil {
		return handleErr
	}
	return c.JSON(200, out)
}

func (s *Server) HandleComAtprotoAdminGetModerationReport(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleComAtprotoAdminGetModerationReport")
	defer span.End()

	id, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		return err
	}
	var out *atproto.AdminDefs_ReportViewDetail
	var handleErr error
	// func (s *Server) handleComAtprotoAdminGetModerationReport(ctx context.Context,id int) (*atproto.AdminDefs_ReportViewDetail, error)
	out, handleErr = s.handleComAtprotoAdminGetModerationReport(ctx, id)
	if handleErr != nil {
		return handleErr
	}
	return c.JSON(200, out)
}

func (s *Server) HandleComAtprotoAdminGetModerationReports(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleComAtprotoAdminGetModerationReports")
	defer span.End()
	before := c.QueryParam("before")

	var limit int
	if p := c.QueryParam("limit"); p != "" {
		var err error
		limit, err = strconv.Atoi(p)
		if err != nil {
			return err
		}
	} else {
		limit = 50
	}

	var resolved *bool
	if p := c.QueryParam("resolved"); p != "" {
		resolved_val, err := strconv.ParseBool(p)
		if err != nil {
			return err
		}
		resolved = &resolved_val
	}
	subject := c.QueryParam("subject")
	var out *atproto.AdminGetModerationReports_Output
	var handleErr error
	// func (s *Server) handleComAtprotoAdminGetModerationReports(ctx context.Context,before string,limit int,resolved *bool,subject string) (*atproto.AdminGetModerationReports_Output, error)
	out, handleErr = s.handleComAtprotoAdminGetModerationReports(ctx, before, limit, resolved, subject)
	if handleErr != nil {
		return handleErr
	}
	return c.JSON(200, out)
}

func (s *Server) HandleComAtprotoAdminResolveModerationReports(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleComAtprotoAdminResolveModerationReports")
	defer span.End()

	var body atproto.AdminResolveModerationReports_Input
	if err := c.Bind(&body); err != nil {
		return err
	}
	var out *atproto.AdminDefs_ActionView
	var handleErr error
	// func (s *Server) handleComAtprotoAdminResolveModerationReports(ctx context.Context,body *atproto.AdminResolveModerationReports_Input) (*atproto.AdminDefs_ActionView, error)
	out, handleErr = s.handleComAtprotoAdminResolveModerationReports(ctx, &body)
	if handleErr != nil {
		return handleErr
	}
	return c.JSON(200, out)
}

func (s *Server) HandleComAtprotoAdminReverseModerationAction(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleComAtprotoAdminReverseModerationAction")
	defer span.End()

	var body atproto.AdminReverseModerationAction_Input
	if err := c.Bind(&body); err != nil {
		return err
	}
	var out *atproto.AdminDefs_ActionView
	var handleErr error
	// func (s *Server) handleComAtprotoAdminReverseModerationAction(ctx context.Context,body *atproto.AdminReverseModerationAction_Input) (*atproto.AdminDefs_ActionView, error)
	out, handleErr = s.handleComAtprotoAdminReverseModerationAction(ctx, &body)
	if handleErr != nil {
		return handleErr
	}
	return c.JSON(200, out)
}

func (s *Server) HandleComAtprotoAdminTakeModerationAction(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleComAtprotoAdminTakeModerationAction")
	defer span.End()

	var body atproto.AdminTakeModerationAction_Input
	if err := c.Bind(&body); err != nil {
		return err
	}
	var out *atproto.AdminDefs_ActionView
	var handleErr error
	// func (s *Server) handleComAtprotoAdminTakeModerationAction(ctx context.Context,body *atproto.AdminTakeModerationAction_Input) (*atproto.AdminDefs_ActionView, error)
	out, handleErr = s.handleComAtprotoAdminTakeModerationAction(ctx, &body)
	if handleErr != nil {
		return handleErr
	}
	return c.JSON(200, out)
}

func (s *Server) HandleComAtprotoReportCreate(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleComAtprotoReportCreate")
	defer span.End()

	var body atproto.ModerationCreateReport_Input
	if err := c.Bind(&body); err != nil {
		return err
	}
	var out *atproto.ModerationCreateReport_Output
	var handleErr error
	// func (s *Server) handleComAtprotoReportCreate(ctx context.Context,body *atproto.ModerationCreateReport_Input) (*atproto.ModerationCreateReport_Output, error)
	out, handleErr = s.handleComAtprotoReportCreate(ctx, &body)
	if handleErr != nil {
		return handleErr
	}
	return c.JSON(200, out)
}

type HealthStatus struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Message string `json:"msg,omitempty"`
}

func (s *Server) HandleHealthCheck(c echo.Context) error {
	if err := s.db.Exec("SELECT 1").Error; err != nil {
		log.Errorf("healthcheck can't connect to database: %v", err)
		return c.JSON(500, HealthStatus{Status: "error", Version: version.Version, Message: "can't connect to database"})
	} else {
		return c.JSON(200, HealthStatus{Status: "ok", Version: version.Version})
	}
}
