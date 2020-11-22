;;; for gauche

(define-module parsertest
  (export pattern pattern/skip
          choice
          many-fold many-fold/skip
          many+fold many+fold/skip
          many many/skip
          many+ many+/skip
          cat cat+
          option
          matches matches/skip
          join join/skip
          join+ join+/skip
          chain chain/skip
          chain-list chain-list/skip
          eos
          define-parser
          parse
          set-skip-parser!))

(select-module parsertest)

(define *unmatch-i* -1)
(define *unmatch-e* '())
(define *skip-parser* #f)

(define-method pattern ((x <string>))
  (lambda (s i)
    (let ((n (string-length x)))
      (if (and (<= (+ i n) (string-length s))
               (string=? x (substring s i (+ i n))))
        (values (substring s i (+ i n)) (+ i n))
        (values #f #f)))))

(define-method pattern ((x <regexp>))
  (lambda (s i)
    (let ((m (rxmatch x (string-copy s i))))
      (if (and m (= 0 (rxmatch-start m)))
        (values (rxmatch-substring m) (+ i (rxmatch-end m)))
        (values #f #f)))))

(define-method pattern (x) x)

(define (pattern/skip x)
  (lambda (s i)
    (receive (drop i) (*skip-parser* s i)
      ((pattern x) s i))))

(define (choice . ps)
  (lambda (s i)
    (let rec ((ps ps))
      (if (pair? ps)
        (receive (v i) ((pattern (car ps)) s i)
          (if i (values v i) (rec (cdr ps))))
        (values #f #f)))))

(define (option p fail)
  (lambda (s i)
    (receive (v ni) ((pattern p) s i)
      (if ni 
        (values v ni)
        (values fail i)))))

(define-syntax many+fold 
  (syntax-rules ()
    ((_ body ...) (_many+fold pattern body ...))))

(define-syntax many+fold/skip
  (syntax-rules ()
    ((_ body ...) (_many+fold pattern/skip body ...))))

(define (_many+fold pt p f init)
  (lambda (s i)
    (receive (v i) ((pt p) s i)
      (if i
        (let rec ((init (f v init)) (i i))
          (receive (v ni) ((pt p) s i)
            (if ni
              (rec (f v init) ni)
              (values init i))))
        (values #f #f)))))

(define (many-fold p f init)
  (option (many+fold p f init) init))

(define (many-fold/skip p f init)
  (option (many+fold/skip p f init) init))

(define (many p)
  (many-fold p (lambda (v init) (append init (list v))) '()))

(define (many/skip p)
  (many-fold/skip p (lambda (v init) (append init (list v))) '()))

(define (many+ p)
  (many+fold p (lambda (v init) (append init (list v))) '()))

(define (many+/skip p)
  (many+fold/skip p (lambda (v init) (append init (list v))) '()))

(define (cat p)
  (many-fold p (lambda (v init) (string-append init v)) ""))

(define (cat+ p)
  (many+fold p (lambda (v init) (string-append init v)) ""))

(define-syntax matches
  (syntax-rules ()
    ((_ body ...) (_matches pattern body ...)))) 

(define-syntax matches/skip
  (syntax-rules ()
    ((_ body ...) (_matches pattern/skip body ...)))) 

(define-syntax _matches
  (syntax-rules (<- =>)
    ((_ pt v <- p => r)
      (lambda (s i)
        (receive (v i) ((pt p) s i)
          (if i (values r i) (values #f #f)))))
    ((_ pt v <- p rest ...)
      (lambda (s i)
        (receive (v i) ((pt p) s i)
          (if i ((_matches pt rest ...) s i) (values #f #f)))))
    ((_ pt p => r)
      (lambda (s i)
        (receive (v i) ((pt p) s i)
          (if i (values r i) (values #f #f)))))
    ((_ pt p rest ...)
      (lambda (s i)
        (receive (v i) ((pt p) s i)
          (if i ((_matches pt rest ...) s i) (values #f #f)))))))

(define (join+ p separator)
  (matches
    f <- p
    r <- (many (matches s <- separator e <- p => e))
    => (cons f r)))

(define (join+/skip p separator)
  (matches/skip
    f <- p
    r <- (many (matches/skip s <- separator e <- p => e))
    => (cons f r)))

(define (join p separator)
  (option (join+ p separator) '()))

(define (join/skip p separator)
  (option (join+/skip p separator) '()))

(define (chain p ope)
  (matches
    f <- p
    r <- (many-fold (matches o <- ope e <- p => (cons o e))
                    (lambda (v init) ((car v) init (cdr v))) f)
    => r))

(define (chain/skip p ope)
  (matches/skip
    f <- p
    r <- (many-fold (matches/skip o <- ope e <- p => (cons o e))
                    (lambda (v init) ((car v) init (cdr v))) f)
    => r))

(define (chain-list p ope)
  (matches
    f <- p
    r <- (many-fold (matches o <- ope e <- p => (cons o e))
                    (lambda (v init) (list (car v) init (cdr v))) f)
    => r))

(define (chain-list/skip p ope)
  (matches/skip
    f <- p
    r <- (many-fold (matches/skip o <- ope e <- p => (cons o e))
                    (lambda (v init) (list (car v) init (cdr v))) f)
    => r))

(define eos
  (lambda (s i)
    (if (<= (string-length s) i)
      (values #f i)
      (values #f #f))))

(define (unmatch v i expect)
  (cond ((< *unmatch-i* i)
          (set! *unmatch-e* (list expect))
          (set! *unmatch-i* i))
        ((= *unmatch-i* i)
          (or (member expect *unmatch-e*)
              (set! *unmatch-e* (cons expect *unmatch-e*)))))
  (values v #f))

(define-syntax define-parser
  (syntax-rules ()
    ((_ name p)
      (define name
        (lambda (s i)
          (receive (v ni) ((pattern p) s i)
            (if ni
              (values v ni)
              (unmatch v i (symbol->string 'name)))))))))

(define (line-and-column s i)
  (let* ((i (min i (string-length s)))
         (ls (string-split (string-copy s 0 i) #/\r\n|\r|\n/))
         (line (length ls))
         (column (string-length (list-ref ls (- line 1)))))
    (values line column)))

(define (error-message s)
  (receive (line column) (line-and-column s *unmatch-i*)
    (format "PARSE FAIL: line:~d column:~d expect-one-of:~s"
            line column *unmatch-e*)))

(define (parse p s)
  (set! *unmatch-i* -1)
  (set! *unmatch-e* '())
  (receive (v i) (p s 0)
    (if i v (error (error-message s)))))

(define (set-skip-parser! p)
  (set! *skip-parser* (pattern p)))

(set-skip-parser! #/[ \t\r\n]*/)



;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;;; infix calculator

; (use parsertest)
;
; (define-parser integer
;   (matches
;     s <- (option "-" "")
;     d <- #/[1-9][0-9]+|[0-9]/
;     => (string->number (string-append s d))))
; 
; (define-parser plus  (matches "+" => '+))
; (define-parser minus (matches "-" => '-))
; (define-parser mul   (matches "*" => '*))
; (define-parser div   (matches "/" => '/))
; 
; (define-parser expr
;   (chain-list/skip term (choice plus minus)))
; 
; (define-parser term
;   (chain-list/skip factor (choice mul div)))
; 
; (define-parser factor
;   (choice
;     (matches/skip i <- integer => i)
;     (matches/skip "(" e <- expr ")" => e)))
; 
; (define-parser calc
;    (matches/skip
;      e <- expr
;      z <- eos
;      => e))

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;;; result

; gosh> (parse calc "1 + 5/3 * (8 + (9 - -4)) / (7*7 + 6) + 2")
; (+ (+ 1 (/ (* (/ 5 3) (+ 8 (- 9 -4))) (+ (* 7 7) 6))) 2)
; gosh> 


;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;;; json parser

; (use parsertest)
; 
; (define-parser json
;   (choice
;     (matches "true"  => 'true)
;     (matches "false" => 'false)
;     (matches "null"  => 'null)
;     json-number
;     json-string
;     json-object
;     json-array))
; 
; (define-parser negative-sign "-")
; (define-parser decimal-part  #/[1-9][0-9]+|[0-9]/)
; (define-parser fraction-part #/\.[0-9]+/)
; (define-parser exponent-part #/[eE][-+]?[0-9]+/)
; 
; (define-parser json-number
;   (matches
;     s <- (option negative-sign "")
;     d <- decimal-part
;     f <- (option fraction-part "")
;     e <- (option exponent-part "")
;     => (string->number (string-append s d f e))))
; 
; (define-parser unicode-escape
;   (matches
;     e <- "\\u"
;     u <- #/[0-9a-fA-F]{4}/
;     => (x->string (ucs->char (string->number u 16)))))
; 
; (define-parser escape-sequence
;   (matches
;     e <- "\\"
;     c <- #/["\\\/bfnrt]/
;     => (cond ((string=? c "b") "\x08")
;              ((string=? c "f") "\x0c")
;              ((string=? c "n") "\n")
;              ((string=? c "r") "\r")
;              ((string=? c "t") "\t")
;              (else (x->string c)))))
; 
; (define-parser string-char #/[^"\\\x00-\x1F]/)
; 
; (define-parser json-string
;   (matches
;     sdq <- "\""
;     str <- (cat (choice unicode-escape escape-sequence string-char))
;     edq <- "\""
;     => str))
; 
; (define-parser colon ":")
; (define-parser comma ",")
; 
; (define-parser json-object-member
;   (matches/skip
;     key <- json-string
;     col <- colon
;     val <- json
;     => (cons key val)))
; 
; (define-parser json-object-members
;   (join/skip json-object-member comma))
; 
; (define-parser json-object
;   (matches/skip
;     lcb <- "{"
;     mem <- json-object-members
;     rcb <- "}"
;     => mem))
; 
; (define-parser json-array-elements
;   (join/skip json comma))
; 
; (define-parser json-array
;   (matches/skip
;     lbr <- "["
;     elm <- json-array-elements
;     rbr <- "]"
;     => (list->vector elm)))
; 
; (define-parser json-parser
;   (matches/skip
;     j <- json
;     end <- eos
;     => j))
; 
; (use file.util)
; (use gauche.time)
; 
; (define (main args)
;   (for-each (lambda (file)
;               (let ((s (file->string file)))
;                 ;(time (format #t "~s" (parse json-parser s)))))
;                 (format #t "~s" (parse json-parser s))))
;             (cdr args))
;   0)

