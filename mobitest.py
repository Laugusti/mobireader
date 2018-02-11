import struct, datetime

class PalmDOCAttributes(object):

    class Attr(object):

        def __init__(self, name, field, val):
            self.name = name
            self.val = val & field

        def __str__(self):
            return '%s: %s'%(self.name, bool(self.val))

    def __init__(self, raw):
        self.val = struct.unpack(b'<H', raw)[0]
        self.attributes = []
        for name, field in [('Read Only', 0x02), ('Dirty AppInfoArea', 0x04),
                ('Backup this database', 0x08),
                ('Okay to install newer over existing copy, if present on PalmPilot', 0x10),
                ('Force the PalmPilot to reset after this database is installed', 0x12),
                ('Don\'t allow copy of file to be beamed to other Pilot',
                    0x14)]:
            self.attributes.append(PalmDOCAttributes.Attr(name, field,
                self.val))

    def __str__(self):
        attrs = '\n\t'.join([str(x) for x in self.attributes])
        return 'PalmDOC Attributes: %s\n\t%s'%(bin(self.val), attrs)


class PalmDB(object):

    def __init__(self, raw):
        self.raw = raw

        if self.raw.startswith(b'TPZ'):
            raise ValueError('This is a Topaz file')

        self.name     = self.raw[:32].replace(b'\x00', b'')
        self.attributes = PalmDOCAttributes(self.raw[32:34])
        self.version = struct.unpack(b'>H', self.raw[34:36])[0]

        palm_epoch = datetime.datetime(1904, 1, 1, tzinfo=None)
        self.creation_date_raw = struct.unpack(b'>I', self.raw[36:40])[0]
        self.creation_date = (palm_epoch +
                datetime.timedelta(seconds=self.creation_date_raw))
        self.modification_date_raw = struct.unpack(b'>I', self.raw[40:44])[0]
        self.modification_date = (palm_epoch +
                datetime.timedelta(seconds=self.modification_date_raw))
        self.last_backup_date_raw = struct.unpack(b'>I', self.raw[44:48])[0]
        self.last_backup_date = (palm_epoch +
                datetime.timedelta(seconds=self.last_backup_date_raw))
        self.modification_number = struct.unpack(b'>I', self.raw[48:52])[0]
        self.app_info_id = self.raw[52:56]
        self.sort_info_id = self.raw[56:60]
        self.type = self.raw[60:64]
        self.creator = self.raw[64:68]
        self.ident = self.type + self.creator
        if self.ident not in (b'BOOKMOBI', b'TEXTREAD'):
            raise ValueError('Unknown book ident: %r'%self.ident)
        self.last_record_uid, = struct.unpack(b'>I', self.raw[68:72])
        self.next_rec_list_id = self.raw[72:76]

        self.number_of_records, = struct.unpack(b'>H', self.raw[76:78])

    def __str__(self):
        ans = ['*'*20 + ' PalmDB Header '+ '*'*20]
        ans.append('Name: %r'%self.name)
        ans.append(str(self.attributes))
        ans.append('Version: %s'%self.version)
        ans.append('Creation date: %s (%s)'%(self.creation_date.isoformat(),
            self.creation_date_raw))
        ans.append('Modification date: %s (%s)'%(self.modification_date.isoformat(),
            self.modification_date_raw))
        ans.append('Backup date: %s (%s)'%(self.last_backup_date.isoformat(),
            self.last_backup_date_raw))
        ans.append('Modification number: %s'%self.modification_number)
        ans.append('App Info ID: %r'%self.app_info_id)
        ans.append('Sort Info ID: %r'%self.sort_info_id)
        ans.append('Type: %r'%self.type)
        ans.append('Creator: %r'%self.creator)
        ans.append('Last record UID +1: %r'%self.last_record_uid)
        ans.append('Next record list id: %r'%self.next_rec_list_id)
        ans.append('Number of records: %s'%self.number_of_records)

        return '\n'.join(ans)

def read_file(raw):
    return PalmDB(raw[:78])

stream = open('data.mobi', 'rb')
raw=stream.read()
palmdb = read_file(raw)
for i in xrange(palmdb.number_of_records):
            pos = 78 + i * 8
            offset, a1, a2, a3, a4 = struct.unpack(b'>LBBBB', raw[pos:pos+8])
            flags, val = a1, a2 << 16 | a3 << 8 | a4
            print(offset, flags, val)
