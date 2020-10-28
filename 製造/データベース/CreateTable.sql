create table store(
    scode char(6),
    sname varchar(255),
    address varchar(255),
    tel varchar(255),
    url varchar(255),
    stime varchar(255),
    capacity int,
    primary key (scode));

create table menu(
    mcode char(4),
    mname varchar(255),
    price int,
    detail varchar(255),
    primary key (mcode));

create table order(
    ocode char(8) not null,
    day date,
    scode char(6),
    tabban int,
    flg boolean,
    primary key(ocode),
    foreign key(scode) references store(scode));

create table list(
    lcode int auto_increment,
    mcode char(4),
    quantity int,
    ocode char(8),
    primary key(lcode),
    foreign key(ocode) references order(ocode));
